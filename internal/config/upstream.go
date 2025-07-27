package config

import (
	"fmt"

	"github.com/buglloc/sly64/v2/internal/config/configpb"
	"github.com/buglloc/sly64/v2/internal/upstream"
)

func NewUpstream(cfg *configpb.Upstream) (upstream.Upstream, error) {
	switch spec := cfg.Kind.(type) {
	case *configpb.Upstream_Udp:
		return upstream.NewPlain(
			upstream.WithAddr(spec.Udp.Addr, upstream.NetUDP),
			upstream.WithDialTimeout(spec.Udp.DialTimeout.AsDuration()),
			upstream.WithTimeout(spec.Udp.Timeout.AsDuration()),
			upstream.WithIface(spec.Udp.Iface),
		)

	case *configpb.Upstream_Tcp:
		return upstream.NewPlain(
			upstream.WithAddr(spec.Tcp.Addr, upstream.NetTCP),
			upstream.WithDialTimeout(spec.Tcp.DialTimeout.AsDuration()),
			upstream.WithTimeout(spec.Tcp.Timeout.AsDuration()),
			upstream.WithIface(spec.Tcp.Iface),
		)

	case *configpb.Upstream_Dot:
		tlsCfg, err := NewTLSConfig(spec.Dot.Tls)
		if err != nil {
			return nil, fmt.Errorf("create TLS config: %w", err)
		}

		return upstream.NewDoT(
			upstream.WithAddr(spec.Dot.Addr, upstream.NetTCPTLS),
			upstream.WithDialTimeout(spec.Dot.DialTimeout.AsDuration()),
			upstream.WithTimeout(spec.Dot.Timeout.AsDuration()),
			upstream.WithIface(spec.Dot.Iface),
			upstream.WithTLSConfig(tlsCfg),
		)

	default:
		return nil, fmt.Errorf("unsupported upstream: %T", spec)
	}
}

func upstreamPatcher(cfg *configpb.Config, path string) error {
	for _, route := range cfg.Route {
		for _, u := range route.Upstream {
			dotU, ok := u.Kind.(*configpb.Upstream_Dot)
			if !ok {
				continue
			}

			patchTLSConfig(dotU.Dot.Tls, path)
		}
	}

	return nil
}
