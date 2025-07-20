package upstream

import "time"

type PlainOption interface {
	isPlainOption()
	Option
	DialOption
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
	addr string
}

func WithAddr(addr string) Option {
	return addrOpt{
		addr: addr,
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
	DialOption
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
