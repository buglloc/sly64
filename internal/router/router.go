package router

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

const (
	buildTrieTimeout = 30 * time.Second
)

type Router struct {
	routes []*Route
	mu     sync.Mutex
	trie   *RouteTrie
}

func NewRouter(routes ...*Route) (*Router, error) {
	r := &Router{
		routes: routes,
	}

	trie, err := r.buildTrie()
	if err != nil {
		return nil, fmt.Errorf("build route trie: %w", err)
	}

	r.trie = trie
	r.subscribe()
	return r, nil
}

func (r *Router) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	if len(req.Question) != 1 {
		return nil, fmt.Errorf("unexpected questions: expected=1 got=%d", len(req.Question))
	}

	r.mu.Lock()
	route := r.trie.Find(req.Question[0].Name)
	r.mu.Unlock()

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

func (r *Router) subscribe() {
	for _, route := range r.routes {
		route.Subscribe(r.onRouteChange)
	}
}

func (r *Router) onRouteChange(route string) {
	log.Info().Str("route", route).Msg("new changes: rebuild trie")

	trie, err := r.buildTrie()
	if err != nil {
		log.Error().Err(err).Msg("failed to build new trie")
		return
	}

	r.mu.Lock()
	r.trie = trie
	r.mu.Unlock()
}

func (r *Router) buildTrie() (*RouteTrie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), buildTrieTimeout)
	defer cancel()

	seen := make(map[string]struct{})
	trie := NewRouteTrie()
	for i, route := range r.routes {
		domains, err := route.Domains(ctx)
		if err != nil {
			return nil, fmt.Errorf("get domains from route %s[%d]: %w", route.Name(), i, err)
		}

		for _, domain := range domains {
			if _, dup := seen[domain]; dup {
				return nil, fmt.Errorf("dupplicate domain in route %s[%d]: %s", route.Name(), i, domain)
			}
			seen[domain] = struct{}{}

			trie.Insert(domain, route)
		}
	}

	return trie, nil
}
