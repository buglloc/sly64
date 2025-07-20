package listener

import (
	"github.com/buglloc/sly64/v2/internal/router"
	"github.com/buglloc/sly64/v2/internal/syncutil"
)

type Option func(*Server)

func WithMaxRequests(max int) Option {
	return func(s *Server) {
		if max <= 0 {
			return
		}

		s.sema = syncutil.NewLeakySemaphore(max)
	}
}

func WithListener(cfg ListenerCfg) Option {
	return func(s *Server) {
		s.servers = append(s.servers, newDNSServer(cfg))
	}
}

func WithRouter(r *router.Router) Option {
	return func(s *Server) {
		s.router = r
	}
}
