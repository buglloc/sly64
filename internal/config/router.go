package config

import (
	"fmt"

	"github.com/buglloc/sly64/v2/internal/config/configpb"
	"github.com/buglloc/sly64/v2/internal/dns64"
	"github.com/buglloc/sly64/v2/internal/router"
	"github.com/buglloc/sly64/v2/internal/upstream"
)

func NewSource(cfg *configpb.Source) (router.Source, error) {
	switch spec := cfg.Kind.(type) {
	case *configpb.Source_Static:
		return router.NewStaticSource(
			spec.Static.Domain,
		)

	case *configpb.Source_File:
		return router.NewFileSource(
			spec.File.Path,
			router.WithFileSourceReloadInterval(spec.File.ReloadInterval.AsDuration()),
		)

	default:
		return nil, fmt.Errorf("unsupported source: %T", spec)
	}
}

func NewRoute(cfg *configpb.Route) (*router.Route, error) {
	var d64 *dns64.DNS64
	if cfg.Dns64 != nil && len(cfg.Dns64.Prefix) > 0 {
		var err error
		d64, err = dns64.NewDNS64(
			dns64.WithPrefix(cfg.Dns64.Prefix),
		)

		if err != nil {
			return nil, fmt.Errorf("create DNS64: %w", err)
		}
	}

	upstreams := make([]upstream.Upstream, len(cfg.Upstream))
	for i, uCfg := range cfg.Upstream {
		var err error
		upstreams[i], err = NewUpstream(uCfg)
		if err != nil {
			return nil, fmt.Errorf("create upstream [%d]: %w", i, err)
		}
	}

	sources := make([]router.Source, len(cfg.Source))
	for i, sCfg := range cfg.Source {
		var err error
		sources[i], err = NewSource(sCfg)
		if err != nil {
			return nil, fmt.Errorf("create source [%d]: %w", i, err)
		}
	}

	var cacheCfg router.CacheCfg
	if cfg.Cache != nil {
		cacheCfg = router.CacheCfg{
			Size:   int(cfg.Cache.MaxItems),
			MinTTL: cfg.Cache.MinTtl,
			MaxTTL: cfg.Cache.MaxTtl,
		}
	}

	return router.NewRoute(
		router.WithRouteName(cfg.Name),
		router.WithRouteSource(sources...),
		router.WithRouteFinalize(cfg.Finalize),
		router.WithRouteUpstreams(upstreams...),
		router.WithRouteDNS64(d64),
		router.WithRouteCache(cacheCfg),
	), nil
}

func (r *Runtime) NewRouter() (*router.Router, error) {
	routes := make([]*router.Route, len(r.Config.Route))
	for i, rCfg := range r.Config.Route {
		var err error
		routes[i], err = NewRoute(rCfg)
		if err != nil {
			return nil, fmt.Errorf("create route [%d]: %w", i, err)
		}
	}

	return router.NewRouter(routes...)
}

func routerPatcher(cfg *configpb.Config, cfgPath string) error {
	for _, route := range cfg.Route {
		for _, source := range route.Source {
			fileSource, ok := source.Kind.(*configpb.Source_File)
			if !ok {
				continue
			}

			fileSource.File.Path = absPath(cfgPath, fileSource.File.Path)
		}
	}

	return nil
}
