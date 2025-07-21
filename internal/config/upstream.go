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
			upstream.WithPlainAddr(spec.Udp.Addr, upstream.NetUDP),
			upstream.WithDialTimeout(spec.Udp.DialTimeout.AsDuration()),
			upstream.WithTimeout(spec.Udp.Timeout.AsDuration()),
			upstream.WithIface(spec.Udp.Iface),
		)

	case *configpb.Upstream_Tcp:
		return upstream.NewPlain(
			upstream.WithPlainAddr(spec.Tcp.Addr, upstream.NetTCP),
			upstream.WithDialTimeout(spec.Tcp.DialTimeout.AsDuration()),
			upstream.WithTimeout(spec.Tcp.Timeout.AsDuration()),
			upstream.WithIface(spec.Tcp.Iface),
		)

	case *configpb.Upstream_Dot:
		return upstream.NewDoT(
			upstream.WithDoTAddr(spec.Dot.Addr, spec.Dot.ServerName),
			upstream.WithDialTimeout(spec.Dot.DialTimeout.AsDuration()),
			upstream.WithTimeout(spec.Dot.Timeout.AsDuration()),
			upstream.WithIface(spec.Dot.Iface),
		)

	default:
		return nil, fmt.Errorf("unsupported upstream: %T", spec)
	}
}
