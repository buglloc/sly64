package upstream

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

const (
	DefaultPlainNetwork = "udp"
	DefaultPlainAddr    = "1.1.1.1:53"
	DefaultPlainTimeout = 2 * time.Second
	DefaultPlainPort    = 53
)

var _ Upstream = (*Plain)(nil)

type plainAddr struct {
	net  string
	addr string
}

func (a plainAddr) String() string {
	return fmt.Sprintf("%s://%s", a.net, a.addr)
}

type Plain struct {
	addr        plainAddr
	canFallback bool
	dialer      *net.Dialer
	timeout     time.Duration
}

func NewPlain(opts ...Option) (*Plain, error) {
	p := &Plain{
		addr: plainAddr{
			net:  DefaultPlainNetwork,
			addr: DefaultPlainAddr,
		},
		canFallback: true,
		timeout:     DefaultPlainTimeout,
	}

	var dialOpts []DialOption
	for _, opt := range opts {
		if o, ok := opt.(DialOption); ok {
			dialOpts = append(dialOpts, o)
			continue
		}

		switch o := opt.(type) {
		case addrOpt:
			addr, err := p.parseAddr(o.addr, o.network)
			if err != nil {
				return nil, fmt.Errorf("invalid upstream addr %q: %w", o.addr, err)
			}
			p.addr = addr

		case timeoutOpt:
			p.timeout = o.timeout

		case nopOpt:
			//pass

		default:
			return nil, fmt.Errorf("unsupported option: %T", o)
		}
	}

	p.dialer = NewDialer(dialOpts...)
	return p, nil
}

func (p *Plain) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	rsp, err := p.exchange(ctx, req, p.addr.net, p.addr.addr)
	if err == nil {
		return rsp, nil
	}

	switch p.addr.net {
	case NetUDP, NetUDP4, NetUDP6:
		return nil, err
	}

	switch {
	case !p.canFallback:
		return nil, err

	case errors.Is(err, ErrMalformedRsp):
		log.Ctx(ctx).
			Info().
			Str("source", "upstream_plain").
			Stringer("addr", p.addr).
			Err(err).
			Msg("plain response is malformed, using tcp")

		return p.exchange(ctx, req, switchNetwork(p.addr.net), p.addr.addr)

	case rsp.Truncated:
		log.Ctx(ctx).
			Info().
			Str("source", "upstream_plain").
			Stringer("addr", p.addr).
			Err(err).
			Msg("plain response is truncated, using tcp")

		return p.exchange(ctx, req, switchNetwork(p.addr.net), p.addr.addr)

	default:
		return nil, err
	}
}

func (p *Plain) exchange(ctx context.Context, req *dns.Msg, network string, addr string) (*dns.Msg, error) {
	dnsc := &dns.Client{
		Net:     network,
		Timeout: p.timeout,
		Dialer:  p.dialer,
	}

	rsp, _, err := dnsc.ExchangeContext(ctx, req, addr)
	if err != nil {
		return nil, fmt.Errorf("exchange with %s: %w", p.addr, err)
	}

	return rsp, validateResponse(req, rsp)
}

func (p *Plain) Address() string {
	return p.addr.String()
}

func (p *Plain) Close() error {
	return nil
}

func (p *Plain) parseAddr(addr string, network string) (plainAddr, error) {
	switch network {
	case NetUDP, NetUDP4, NetUDP6, NetTCP, NetTCP4, NetTCP6:
		// pass

	default:
		return plainAddr{}, fmt.Errorf("unsupported network %s", network)
	}

	if _, portStr, err := net.SplitHostPort(addr); err == nil {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return plainAddr{}, fmt.Errorf("invalid port %s: %w", portStr, err)
		}

		if port < 0 || port > 65535 {
			return plainAddr{}, fmt.Errorf("invalid port %s: out of range", portStr)
		}

		return plainAddr{
			net:  network,
			addr: addr,
		}, nil
	}

	return plainAddr{
		net:  network,
		addr: fmt.Sprintf("%s:%d", addr, DefaultPlainPort),
	}, nil
}
