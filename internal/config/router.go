package config

import (
	"fmt"

	"github.com/buglloc/sly64/v2/internal/dns64"
	"github.com/buglloc/sly64/v2/internal/router"
	"github.com/buglloc/sly64/v2/internal/upstream"
)

type Route struct {
	Name      string     `yaml:"name"`
	Finalize  bool       `yaml:"finalize"`
	DNS64     DNS64      `yaml:"dns64"`
	Cache     Cache      `yaml:"cache"`
	Upstreams []Upstream `yaml:"upstreams"`
	Domains   []string   `yaml:"domains"`
}

type Cache struct {
	Size   int    `yaml:"size"`
	MinTTL uint32 `yaml:"min_ttl"`
	MaxTTL uint32 `yaml:"max_ttl"`
}

type DNS64 struct {
	Prefix string `yaml:"prefix"`
}

func NewRoute(cfg Route) (*router.Route, error) {
	var d64 *dns64.DNS64
	if len(cfg.DNS64.Prefix) > 0 {
		var err error
		d64, err = dns64.NewDNS64(
			dns64.WithPrefix(cfg.DNS64.Prefix),
		)

		if err != nil {
			return nil, fmt.Errorf("create DNS64: %w", err)
		}
	}

	upstreams := make([]upstream.Upstream, len(cfg.Upstreams))
	for i, uCfg := range cfg.Upstreams {
		var err error
		upstreams[i], err = NewUpstream(uCfg)
		if err != nil {
			return nil, fmt.Errorf("create upstream [%d]: %w", i, err)
		}
	}

	return router.NewRoute(
		router.WithRouteCache(router.CacheCfg{
			Size:   cfg.Cache.Size,
			MinTTL: cfg.Cache.MinTTL,
			MaxTTL: cfg.Cache.MaxTTL,
		}),
		router.WithRouteName(cfg.Name),
		router.WithRouteDomains(cfg.Domains...),
		router.WithRouteFinalize(cfg.Finalize),
		router.WithRouteUpstreams(upstreams...),
		router.WithRouteDNS64(d64),
	), nil
}

func (r *Runtime) NewRouter() (*router.Router, error) {
	routes := make([]*router.Route, len(r.Config.Routes))
	for i, rCfg := range r.Config.Routes {
		var err error
		routes[i], err = NewRoute(rCfg)
		if err != nil {
			return nil, fmt.Errorf("create route [%d]: %w", i, err)
		}
	}

	return router.NewRouter(routes...)
}
