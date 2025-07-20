package router

import (
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker/v2"

	"github.com/buglloc/sly64/v2/internal/dns64"
	"github.com/buglloc/sly64/v2/internal/upstream"
)

type RouteOption func(r *Route)

func WithRouteName(name string) RouteOption {
	return func(r *Route) {
		r.name = name
	}
}

func WithRouteCache(cfg CacheCfg) RouteOption {
	return func(r *Route) {
		r.cache = NewCache(cfg)
	}
}

func WithRouteUpstreams(upstreams ...upstream.Upstream) RouteOption {
	return func(r *Route) {
		r.upstreams = upstreams
		r.breakers = make([]*gobreaker.CircuitBreaker[*dns.Msg], len(r.upstreams))
		for i, u := range r.upstreams {
			r.breakers[i] = gobreaker.NewCircuitBreaker[*dns.Msg](gobreaker.Settings{
				Name: u.Address(),
				OnStateChange: func(name string, from, to gobreaker.State) {
					log.Warn().
						Str("source", "route_upstream").
						Str("name", name).
						Stringer("form", from).
						Stringer("to", to).
						Msg("upstream circuit breaker state changed")
				},
			})
		}
	}
}

func WithRouteDomains(domains ...string) RouteOption {
	return func(r *Route) {
		if len(domains) == 0 {
			return
		}

		r.domains = domains
	}
}

func WithRouteDNS64(d64 *dns64.DNS64) RouteOption {
	return func(r *Route) {
		r.dns64 = d64
	}
}

func WithRouteFinalize(finalize bool) RouteOption {
	return func(r *Route) {
		r.finalize = finalize
	}
}
