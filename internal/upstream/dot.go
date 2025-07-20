package upstream

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/buglloc/certifi"
	"github.com/miekg/dns"
)

const (
	DefaultDoTNetwork    = "tcp-tls"
	DefaultDoTAddr       = "1.1.1.1:853"
	DefaultDotServerName = "one.one.one.one"
	DefaultDoTTimeout    = 2 * time.Second
	DefaultDoTPort       = 853
)

var _ Upstream = (*DoT)(nil)

type dotAddr struct {
	net        string
	addr       string
	serverName string
}

func (a dotAddr) String() string {
	return fmt.Sprintf("%s://%s [%s]", a.net, a.addr, a.serverName)
}

type DoT struct {
	addr      dotAddr
	dialer    *net.Dialer
	timeout   time.Duration
	tlsConfig *tls.Config
}

func NewDoT(opts ...Option) (*DoT, error) {
	d := &DoT{
		addr: dotAddr{
			net:        DefaultDoTNetwork,
			addr:       DefaultDoTAddr,
			serverName: DefaultDotServerName,
		},
		timeout: DefaultPlainTimeout,
		tlsConfig: &tls.Config{
			RootCAs:            certifi.NewCertPool(),
			ClientSessionCache: tls.NewLRUClientSessionCache(0),
			MinVersion:         tls.VersionTLS12,
		},
	}

	var dialOpts []DialOption
	for _, opt := range opts {
		if o, ok := opt.(DialOption); ok {
			dialOpts = append(dialOpts, o)
			continue
		}

		switch o := opt.(type) {
		case dotAddrOpt:
			addr, err := d.parseAddr(o.addr, o.serverName)
			if err != nil {
				return nil, fmt.Errorf("invalid upstream addr %q: %w", o.addr, err)
			}
			d.addr = addr
			d.tlsConfig.ServerName = addr.serverName

		case timeoutOpt:
			d.timeout = o.timeout

		case nopOpt:
			//pass

		default:
			return nil, fmt.Errorf("unsupported option: %T", o)
		}
	}

	d.dialer = NewDialer(dialOpts...)
	return d, nil
}

func (d *DoT) Exchange(ctx context.Context, req *dns.Msg) (*dns.Msg, error) {
	dnsc := &dns.Client{
		Net:       d.addr.net,
		Timeout:   d.timeout,
		TLSConfig: d.tlsConfig,
		Dialer:    d.dialer,
	}

	rsp, _, err := dnsc.ExchangeContext(ctx, req, d.addr.addr)
	if err != nil {
		return nil, fmt.Errorf("exchange with %s: %w", d.addr, err)
	}

	return rsp, validateResponse(req, rsp)
}

func (d *DoT) Address() string {
	return d.addr.String()
}

func (d *DoT) Close() error {
	return nil
}

func (d *DoT) parseAddr(addr, serverName string) (dotAddr, error) {
	if len(serverName) == 0 {
		return dotAddr{}, errors.New("servername can't be empty so far")
	}

	var uu *url.URL
	switch {
	case strings.Contains(addr, "://"):
		var err error
		uu, err = url.Parse(addr)
		if err != nil {
			return dotAddr{}, fmt.Errorf("parse url: %w", err)
		}
	default:
		uu = &url.URL{
			Scheme: "tcp-tls",
			Host:   addr,
		}
	}

	var network string
	switch uu.Scheme {
	case netTCPTLS:
		network = uu.Scheme
	default:
		return dotAddr{}, fmt.Errorf("unsupported scheme %s", uu.Scheme)
	}

	if _, portStr, err := net.SplitHostPort(uu.Host); err == nil {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return dotAddr{}, fmt.Errorf("invalid port %s: %w", portStr, err)
		}

		if port < 0 || port > 65535 {
			return dotAddr{}, fmt.Errorf("invalid port %s: out of range", portStr)
		}

		return dotAddr{
			net:        network,
			addr:       uu.Host,
			serverName: serverName,
		}, nil
	}

	return dotAddr{
		net:        network,
		addr:       fmt.Sprintf("%s:%d", uu.Host, DefaultPlainPort),
		serverName: serverName,
	}, nil
}
