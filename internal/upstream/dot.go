package upstream

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/buglloc/certifi"
	"github.com/miekg/dns"
)

const (
	DefaultDoTNetwork = "tcp-tls"
	DefaultDoTAddr    = "1.1.1.1:853"
	DefaultDoTTimeout = 2 * time.Second
	DefaultDoTPort    = 853
)

var _ Upstream = (*DoT)(nil)

type dotAddr struct {
	net        string
	addr       string
	serverName string
}

func (a dotAddr) String() string {
	if len(a.serverName) == 0 {
		return fmt.Sprintf("%s://%s", a.net, a.addr)
	}

	return fmt.Sprintf("%s://%s [%s]", a.net, a.addr, a.serverName)
}

type DoT struct {
	addr      dotAddr
	dialer    *net.Dialer
	tlsConfig *tls.Config
	dnsc      *dns.Client
	pool      ConnPool
}

func NewDoT(opts ...Option) (*DoT, error) {
	d := &DoT{
		addr: dotAddr{
			net:  DefaultDoTNetwork,
			addr: DefaultDoTAddr,
		},
		tlsConfig: &tls.Config{
			RootCAs:            certifi.NewCertPool(),
			ClientSessionCache: tls.NewLRUClientSessionCache(0),
		},
		dnsc: &dns.Client{
			Timeout: DefaultDoTTimeout,
		},
	}

	var dialOpts []DialOption
	poolSize := int32(0)
	for _, opt := range opts {
		if o, ok := opt.(DialOption); ok {
			dialOpts = append(dialOpts, o)
			continue
		}

		switch o := opt.(type) {
		case addrOpt:
			addr, err := d.parseAddr(o.addr)
			if err != nil {
				return nil, fmt.Errorf("invalid upstream addr %q: %w", o.addr, err)
			}
			d.addr = addr

		case timeoutOpt:
			d.dnsc.Timeout = o.timeout

		case tlsCfgOpt:
			d.tlsConfig = o.cfg
			d.addr.serverName = o.cfg.ServerName

		case poolOpt:
			poolSize = o.maxItems

		case nopOpt:
			//pass

		default:
			return nil, fmt.Errorf("unsupported option: %T", o)
		}
	}

	d.dialer = NewDialer(dialOpts...)
	dialerFn := func(ctx context.Context) (net.Conn, error) {
		fmt.Println("new conn")
		return tls.DialWithDialer(d.dialer, "tcp", d.addr.addr, d.tlsConfig)
	}

	if poolSize > 0 {
		pool, err := NewPuddlePool(dialerFn, poolSize)
		if err != nil {
			return nil, fmt.Errorf("create puddle pool: %w", err)
		}

		d.pool = pool
	} else {
		d.pool = NewNetPool(dialerFn)
	}

	return d, nil
}

func (d *DoT) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	var errs []error
	for range 2 {
		conn, err := d.pool.Acquire(ctx)
		if err != nil {
			return nil, fmt.Errorf("acquire DoT conn from pool: %w", err)
		}

		rsp, _, err := d.dnsc.ExchangeWithConnContext(ctx, req, &dns.Conn{
			Conn: conn,
		})
		if err != nil {
			// Connection failed, destroy it and try again
			conn.Destroy()
			errs = append(errs, err)
			continue
		}

		// Success, return connection to pool
		conn.Close()
		return rsp, validateResponse(req, rsp)
	}

	return nil, fmt.Errorf("exchange failed after retries: %w", errors.Join(errs...))
}

func (d *DoT) Address() string {
	return d.addr.String()
}

func (d *DoT) Close() error {
	d.pool.Close()
	return nil
}

func (d *DoT) parseAddr(addr string) (dotAddr, error) {
	if _, portStr, err := net.SplitHostPort(addr); err == nil {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return dotAddr{}, fmt.Errorf("invalid port %s: %w", portStr, err)
		}

		if port < 0 || port > 65535 {
			return dotAddr{}, fmt.Errorf("invalid port %s: out of range", portStr)
		}

		return dotAddr{
			net:  "tcp-tls",
			addr: addr,
		}, nil
	}

	return dotAddr{
		net:  "tcp-tls",
		addr: fmt.Sprintf("%s:%d", addr, DefaultDoTPort),
	}, nil
}
