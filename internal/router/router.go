package router

import (
	"context"
	"errors"
	"fmt"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type Router struct {
	routes []*Route
	trie   *RouteTrie
}

func NewRouter(routes ...*Route) (*Router, error) {
	r := &Router{
		routes: routes,
		trie:   NewRouteTrie(),
	}

	seen := make(map[string]struct{})
	for _, route := range r.routes {
		for _, domain := range route.Domains() {
			if _, dup := seen[domain]; dup {
				return nil, fmt.Errorf("dupplicate domain: %s", domain)
			}
			seen[domain] = struct{}{}

			r.trie.Insert(domain, route)
		}
	}

	return r, nil
}

func (r *Router) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	if len(req.Question) != 1 {
		return nil, fmt.Errorf("unexpected questions: expected=1 got=%d", len(req.Question))
	}

	route := r.trie.Find(req.Question[0].Name)
	if route == nil {
		log.Ctx(ctx).Warn().Msgf("no route for name: %s", req.Question[0].Name)
		return &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Rcode: dns.RcodeSuccess,
			},
		}, nil
	}

	return route.Exchange(ctx, req)
}

func (r *Router) Close() error {
	var errs []error
	for _, route := range r.routes {
		if err := route.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close route %s: %w", route.Name(), err))
		}
	}

	return errors.Join(errs...)
}
