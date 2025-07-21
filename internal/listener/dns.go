package listener

import (
	"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type ListenerCfg struct {
	// Address to listen on, ":dns" if empty.
	Addr string
	// if "tcp" or "tcp-tls" (DNS over TLS) it will invoke a TCP listener, otherwise an UDP one
	Net string
	// The net.Conn.SetReadTimeout value for new connections, defaults to 2 * time.Second.
	ReadTimeout time.Duration
	// The net.Conn.SetWriteTimeout value for new connections, defaults to 2 * time.Second.
	WriteTimeout time.Duration
}

func newDNSServer(cfg ListenerCfg) *dns.Server {
	return &dns.Server{
		Addr:         cfg.Addr,
		Net:          cfg.Net,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		NotifyStartedFunc: func() {
			log.Info().Msgf("dns server started: net=%s addr=%s", cfg.Net, cfg.Addr)
		},
	}
}
