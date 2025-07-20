package upstream

import (
	"context"

	"github.com/miekg/dns"
)

type Upstream interface {
	// Exchange sends req to this upstream and returns the response that has
	// been received or an error if something went wrong.  The implementations
	// must not modify req as well as the caller must not modify it until the
	// method returns.  It shouldn't be called after closing.
	Exchange(ctx context.Context, req *dns.Msg) (rsp *dns.Msg, err error)

	Address() (addr string)
	Close() error
}
