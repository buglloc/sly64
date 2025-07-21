package upstream

import (
	"context"
	"fmt"
	"strings"

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

func validateResponse(req, resp *dns.Msg) (err error) {
	if qlen := len(resp.Question); qlen != 1 {
		return fmt.Errorf("%w: only 1 question allowed; got %d", ErrMalformedRsp, qlen)
	}

	reqQ, respQ := req.Question[0], resp.Question[0]
	if reqQ.Qtype != respQ.Qtype {
		return fmt.Errorf("%w: mismatched type %s", ErrMalformedRsp, dns.Type(respQ.Qtype))
	}

	// Compare the names case-insensitively, just like CoreDNS does.
	if !strings.EqualFold(reqQ.Name, respQ.Name) {
		return fmt.Errorf("%w: mismatched name %q", ErrMalformedRsp, respQ.Name)
	}

	return nil
}
