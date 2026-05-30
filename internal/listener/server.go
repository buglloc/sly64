package listener

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/buglloc/sly64/v2/internal/router"
	"github.com/buglloc/sly64/v2/internal/syncutil"
)

const (
	listenerShutdownTimeout = 5 * time.Second
)

var _ dns.Handler = (*Server)(nil)

type Server struct {
	router     *router.Router
	sema       syncutil.Semaphore
	servers    []*dns.Server
	closed     chan struct{}
	ctx        context.Context
	shutdownFn context.CancelFunc
}

func NewServer(opts ...Option) (*Server, error) {
	ctx, shutdownFn := context.WithCancel(context.Background())
	s := &Server{
		sema:       syncutil.NopSemaphore{},
		closed:     make(chan struct{}),
		ctx:        ctx,
		shutdownFn: shutdownFn,
	}

	for _, opt := range opts {
		opt(s)
	}

	if len(s.servers) == 0 {
		return nil, errors.New("no listeners provided")
	}

	for i := range s.servers {
		s.servers[i].Handler = s
	}

	return s, nil
}

func (s *Server) Start() error {
	defer close(s.closed)

	g, ctx := errgroup.WithContext(s.ctx)
	for _, srv := range s.servers {
		srv := srv
		g.Go(func() error {
			err := srv.ListenAndServe()
			if err == nil || errors.Is(err, context.Canceled) {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			default:
			}

			log.Error().Err(err).Msg("listen failed")
			return fmt.Errorf("start server: %w", err)
		})
	}

	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), listenerShutdownTimeout)
		defer cancel()
		for _, srv := range s.servers {
			if err := srv.ShutdownContext(shutdownCtx); err != nil {
				log.Error().Err(err).Msg("listener shutdown failed")
			}
		}

		return nil
	})

	return g.Wait()
}

func (s *Server) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	rsp := new(dns.Msg)
	rsp.SetReply(req)

	if !s.sema.TryAcquire() {
		log.Warn().Msg("too many concurrent requests")
		rsp.SetRcode(req, dns.RcodeServerFailure)
		_ = w.WriteMsg(rsp)
		return
	}
	defer s.sema.Release()

	s.buildReply(req, rsp)
	_ = w.WriteMsg(rsp)
}

func (s *Server) Shutdown(ctx context.Context) error {
	defer func() {
		s.servers = s.servers[:0]
	}()

	s.shutdownFn()

	var errs []error
	for _, srv := range s.servers {
		if err := srv.ShutdownContext(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown listener: %w", err))
		}
	}
	s.servers = s.servers[:0]

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.closed:
		if s.router != nil {
			if err := s.router.Close(); err != nil {
				errs = append(errs, fmt.Errorf("close router: %w", err))
			}
		}

		return errors.Join(errs...)
	}
}

func (s *Server) buildReply(req *dns.Msg, rsp *dns.Msg) {
	if err := s.validateRequest(req, rsp); err != nil {
		log.Warn().Err(err).Stringer("req", req).Msg("invalid request")
		return
	}

	q := req.Question[0]
	ctx := log.With().
		Str("qtype", dns.TypeToString[q.Qtype]).
		Str("qname", q.Name).
		Logger().
		WithContext(s.ctx)

	rrsp, err := s.router.Exchange(ctx, req)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("proxy request failed")
		rsp.SetRcode(req, dns.RcodeServerFailure)
		return
	}

	rsp.MsgHdr = rrsp.MsgHdr
	rsp.Id = req.Id
	rsp.Response = true
	rsp.Opcode = req.Opcode
	rsp.RecursionDesired = req.RecursionDesired
	rsp.CheckingDisabled = req.CheckingDisabled
	rsp.Question = slices.Clone(req.Question)
	rsp.Answer = rrsp.Answer
	rsp.Ns = rrsp.Ns
	rsp.Extra = rrsp.Extra
}

func (s *Server) validateRequest(req *dns.Msg, rsp *dns.Msg) error {
	switch {
	case len(req.Question) != 1:
		rsp.SetRcode(req, dns.RcodeRefused)
		return fmt.Errorf("expected 1 question, got: %d", len(req.Question))

	case req.Question[0].Qtype == dns.TypeANY:
		rsp.SetRcode(req, dns.RcodeNotImplemented)
		return errors.New("refusing dns type any request")

	case s.router == nil:
		rsp.SetRcode(req, dns.RcodeServerFailure)
		return errors.New("can't handle message w/o router")

	default:
		return nil
	}
}
