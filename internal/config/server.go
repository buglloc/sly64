package config

import (
	"fmt"

	"github.com/buglloc/sly64/v2/internal/config/configpb"
	"github.com/buglloc/sly64/v2/internal/listener"
)

func (r *Runtime) NewServer() (*listener.Server, error) {
	router, err := r.NewRouter()
	if err != nil {
		return nil, fmt.Errorf("create router: %w", err)
	}

	opts := []listener.Option{
		listener.WithMaxRequests(int(r.Config.MaxRequests)),
		listener.WithRouter(router),
	}

	for i, cfg := range r.Config.Listener {
		ln, err := parseProtoListenNet(cfg.Net)
		if err != nil {
			return nil, fmt.Errorf("invalid listener net [%d]: %w", i, err)
		}

		opts = append(opts, listener.WithListener(listener.ListenerCfg{
			Addr:         cfg.Addr,
			Net:          ln,
			ReadTimeout:  cfg.ReadTimeout.AsDuration(),
			WriteTimeout: cfg.WriteTimeout.AsDuration(),
		}))
	}

	return listener.NewServer(opts...)
}

func parseProtoListenNet(n configpb.Net) (string, error) {
	switch n {
	case configpb.Net_NET_UDP:
		return "udp", nil

	case configpb.Net_NET_TCP:
		return "tcp", nil

	default:
		return "", fmt.Errorf("unsupported net: %s", n)
	}
}
