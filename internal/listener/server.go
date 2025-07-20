package listener

import (
	"context"
	"errors"
	"fmt"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/buglloc/sly64/v2/internal/router"
	"github.com/buglloc/sly64/v2/internal/syncutil"
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
		g.Go(func() error {
			err := srv.ListenAndServe()
			if err != nil {
				log.Error().Err(err).Msg("listen failed")
				return fmt.Errorf("start server: %w", err)
			}

			return nil
		})
	}

	g.Go(func() error {
		select {
		case <-s.ctx.Done():
		case <-ctx.Done():
		}

		return nil
	})

	return g.Wait()
}

func (s *Server) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	rsp := new(dns.Msg)
	rsp.SetReply(req)

	if err := s.sema.Acquire(s.ctx); err != nil {
		log.Error().Err(err).Msg("acquiring semaphore")
		rsp.SetRcode(req, dns.RcodeServerFailure)
		return
	}
	defer s.sema.Release()

	s.buildReply(req, rsp)
	_ = w.WriteMsg(rsp)
}

func (s *Server) Shutdown(ctx context.Context) error {
	for _, srv := range s.servers {
		srv.ShutdownContext(ctx)
	}
	s.servers = s.servers[:0]

	s.shutdownFn()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.closed:
		return nil
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

	rsp.SetRcode(req, rrsp.Rcode)
	rsp.Answer = rrsp.Answer
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
