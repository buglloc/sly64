package config

import (
	"fmt"
	"time"

	"github.com/buglloc/sly64/v2/internal/listener"
)

type Listener struct {
	// Address to listen on, ":dns" if empty.
	Addr string `yaml:"addr"`
	// if "tcp" or "tcp-tls" (DNS over TLS) it will invoke a TCP listener, otherwise an UDP one
	Net string `yaml:"net"`
	// The net.Conn.SetReadTimeout value for new connections, defaults to 2 * time.Second.
	ReadTimeout time.Duration `yaml:"read_timeout"`
	// The net.Conn.SetWriteTimeout value for new connections, defaults to 2 * time.Second.
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

func (r *Runtime) NewServer() (*listener.Server, error) {
	router, err := r.NewRouter()
	if err != nil {
		return nil, fmt.Errorf("create router: %w", err)
	}

	opts := []listener.Option{
		listener.WithMaxRequests(r.Config.MaxRequests),
		listener.WithRouter(router),
	}

	for _, cfg := range r.Config.Listeners {
		opts = append(opts, listener.WithListener(listener.ListenerCfg{
			Addr:         cfg.Addr,
			Net:          cfg.Net,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		}))
	}

	return listener.NewServer(opts...)
}
