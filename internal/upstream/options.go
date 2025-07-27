package upstream

import (
	"crypto/tls"
	"time"
)

type PlainOption interface {
	isPlainOption()
	Option
}

type DoTOption interface {
	isDoTOption()
	Option
}

type DialOption interface {
	isDialOption()
	Option
}

type Option interface {
	isOption()
}

type nopDialOpt struct {
	DialOption
}

type nopOpt struct {
	Option
}

type addrOpt struct {
	Option
	addr    string
	network string
}

func WithAddr(addr string, network string) Option {
	return addrOpt{
		addr:    addr,
		network: network,
	}
}

type dialTimeoutOpt struct {
	DialOption
	timeout time.Duration
}

func WithDialTimeout(timeout time.Duration) DialOption {
	if timeout == 0 {
		return nopDialOpt{}
	}

	return dialTimeoutOpt{
		timeout: timeout,
	}
}

type timeoutOpt struct {
	Option
	timeout time.Duration
}

func WithTimeout(timeout time.Duration) Option {
	if timeout == 0 {
		return nopOpt{}
	}

	return timeoutOpt{
		timeout: timeout,
	}
}

type ifaceOpt struct {
	DialOption
	iface string
}

func WithIface(iface string) DialOption {
	if len(iface) == 0 {
		return nopDialOpt{}
	}

	return ifaceOpt{
		iface: iface,
	}
}

type tlsCfgOpt struct {
	Option
	cfg *tls.Config
}

func WithTLSConfig(cfg *tls.Config) Option {
	if cfg == nil {
		return nopOpt{}
	}

	if cfg.ClientSessionCache == nil {
		cfg.ClientSessionCache = tls.NewLRUClientSessionCache(0)
	}

	return tlsCfgOpt{
		cfg: cfg,
	}
}
