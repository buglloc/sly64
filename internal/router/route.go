package router

import (
	"context"
	"errors"
	"fmt"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker/v2"

	"github.com/buglloc/sly64/v2/internal/dns64"
	"github.com/buglloc/sly64/v2/internal/upstream"
)

const maxCNAMEDepth = 10

type RouteSubscribeFn func(routeName string)
type Route struct {
	name      string
	upstreams []upstream.Upstream
	sources   []Source
	breakers  []*gobreaker.CircuitBreaker[*dns.Msg]
	finalize  bool
	dns64     *dns64.DNS64
	cache     Cache
}

func NewRoute(opts ...RouteOption) *Route {
	r := &Route{
		cache: &nopCache{},
		sources: []Source{
			&StaticSource{
				domains: []string{
					"*.",
				},
			},
		},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Route) Name() string {
	return r.name
}

func (r *Route) Subscribe(fn RouteSubscribeFn) {
	for _, s := range r.sources {
		s.Watch(func() {
			fn(r.name)
		})
	}
}

func (r *Route) Domains(ctx context.Context) ([]string, error) {
	var errs []error
	var out []string
	for i, s := range r.sources {
		domains, err := s.Domains(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("get domains from source [%d]: %w", i, err))
			continue
		}

		out = append(out, domains...)
	}

	return out, errors.Join(errs...)
}

func (r *Route) Close() error {
	var errs []error
	for _, u := range r.upstreams {
		if err := u.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close upstream %s: %w", u.Address(), err))
		}
	}

	return errors.Join(errs...)
}

func (r *Route) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	if r.dns64 != nil {
		return r.exchangeDNS64(ctx, req)
	}

	return r.exchange(ctx, req)
}

func (r *Route) exchangeDNS64(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	if len(req.Question) != 1 {
		return nil, fmt.Errorf("unexpected questions: expected=1 got=%d", len(req.Question))
	}

	switch req.Question[0].Qtype {
	case dns.TypeSRV, dns.TypeANY, dns.TypeHTTPS, dns.TypeCNAME, dns.TypeA:
		return &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Rcode: dns.RcodeSuccess,
			},
		}, nil

	case dns.TypeAAAA:
		if r.dns64 == nil {
			return nil, errors.New("no DNS64 configured")
		}

		aReq := req.Copy()
		aReq.Question[0].Qtype = dns.TypeA
		rsp, err := r.exchange(ctx, aReq)
		if err != nil {
			return nil, err
		}

		if rsp.Rcode != dns.RcodeSuccess {
			return rsp, nil
		}

		answers := make([]dns.RR, 0, len(rsp.Answer))
		for _, rr := range rsp.Answer {
			aRec, ok := rr.(*dns.A)
			if !ok {
				continue
			}

			v6, err := r.dns64.To6(aRec.A)
			if err != nil {
				continue
			}

			answers = append(answers, &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   req.Question[0].Name,
					Rrtype: dns.TypeAAAA,
					Class:  req.Question[0].Qclass,
					Ttl:    rr.Header().Ttl,
				},
				AAAA: v6,
			})
		}

		rsp.Answer = answers
		return rsp, nil

	default:
		return r.exchange(ctx, req)
	}
}

func (r *Route) exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	if len(req.Question) != 1 {
		return nil, fmt.Errorf("unexpected questions: expected=1 got=%d", len(req.Question))
	}

	if !r.finalize {
		return r.doExchange(ctx, req)
	}

	qtype := req.Question[0].Qtype
	switch qtype {
	case dns.TypeCNAME:
		return &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Rcode: dns.RcodeSuccess,
			},
		}, nil

	case dns.TypeA, dns.TypeAAAA:
		//pass

	default:
		return r.doExchange(ctx, req)
	}

	subReq := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			RecursionDesired: true,
			Opcode:           dns.OpcodeQuery,
		},
		Question: []dns.Question{{
			Qtype:  qtype,
			Qclass: dns.ClassINET,
		}},
	}

	currentName := req.Question[0].Name
	for depth := 0; depth < maxCNAMEDepth; depth++ {
		subReq.Question[0].Name = currentName
		rsp, err := r.doExchange(ctx, subReq)
		if err != nil {
			return nil, err
		}

		if rsp.Rcode != dns.RcodeSuccess {
			return rsp, nil
		}

		// Collect A/AAAA responses
		if rrs := parseAnswers(currentName, qtype, dns.ClassINET, rsp); len(rrs) > 0 {
			rsp.Answer = rrs
			return rsp, nil
		}

		// Find CNAME in the answer section for the current name
		cname := ""
		for _, rr := range rsp.Answer {
			if rr.Header().Rrtype != dns.TypeCNAME {
				continue
			}

			cname = rr.(*dns.CNAME).Target
		}

		if len(cname) > 0 {
			if rrs := parseAnswers(cname, qtype, dns.ClassINET, rsp); len(rrs) > 0 {
				rsp.Answer = rsp.Answer[:0]
				for _, rr := range rrs {
					rr.Header().Name = currentName
					rsp.Answer = append(rsp.Answer, rr)
				}

				return rsp, nil
			}

			currentName = cname
			continue // Follow the CNAME
		}

		return rsp, nil
	}

	return nil, errors.New("CNAME chain too deep or loop detected")
}

func (r *Route) doExchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	return r.cache.Fetch(req, func() (*dns.Msg, error) {
		for i, u := range r.upstreams {
			rsp, err := r.breakers[i].Execute(func() (*dns.Msg, error) {
				return u.Exchange(ctx, req)
			})

			if err != nil {
				if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
					continue
				}

				log.Ctx(ctx).Warn().
					Err(err).
					Str("source", "route").
					Str("route_name", r.name).
					Str("upstream", u.Address()).
					Msg("proxy request failed")
				continue
			}

			return rsp, nil
		}

		return nil, errors.New("no live upstream available")
	})
}

func parseAnswers(qname string, qtype uint16, qclass uint16, rsp *dns.Msg) []dns.RR {
	answers := parseAnswersFromRR(qname, qtype, qclass, rsp.Answer)
	if len(answers) > 0 {
		return answers
	}

	return parseAnswersFromRR(qname, qtype, qclass, rsp.Extra)
}

func parseAnswersFromRR(qname string, qtype uint16, qclass uint16, rrs []dns.RR) []dns.RR {
	answers := make([]dns.RR, 0, len(rrs))
	for _, rr := range rrs {
		header := rr.Header()
		if header.Class != qclass {
			continue
		}

		if header.Name != qname {
			continue
		}

		if header.Rrtype != qtype {
			continue
		}

		answers = append(answers, rr)
	}

	return answers
}
