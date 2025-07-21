package upstream

import "time"

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

type plainAddrOpt struct {
	PlainOption
	addr    string
	network string
}

func WithPlainAddr(addr string, network string) PlainOption {
	return plainAddrOpt{
		addr:    addr,
		network: network,
	}
}

type dotAddrOpt struct {
	DoTOption
	addr       string
	serverName string
}

func WithDoTAddr(addr, serverName string) DoTOption {
	return dotAddrOpt{
		addr:       addr,
		serverName: serverName,
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
