package upstream

import (
	"context"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	DefaultDialTimeout = 500 * time.Millisecond

	NetUDP  = "udp"
	NetUDP4 = "udp4"
	NetUDP6 = "udp6"

	NetTCP  = "tcp"
	NetTCP4 = "tcp4"
	NetTCP6 = "tcp6"

	NetTCPTLS = "tcp-tls"
)

func NewDialer(opts ...DialOption) *net.Dialer {
	d := &net.Dialer{}

	for _, opt := range opts {
		switch o := opt.(type) {
		case dialTimeoutOpt:
			d.Timeout = o.timeout

		case ifaceOpt:
			d.ControlContext = newBindToDeviceControl(o.iface)

		case nopDialOpt:

		default:
			log.Error().
				Str("source", "dialer").
				Type("option", o).
				Msg("skip unsupported option")
		}
	}

	return d
}

func newBindToDeviceControl(dev string) func(ctx context.Context, network, address string, c syscall.RawConn) error {
	return func(ctx context.Context, network, address string, c syscall.RawConn) error {
		var cerr error
		err := c.Control(func(fd uintptr) {
			cerr = syscall.BindToDevice(int(fd), dev)
		})

		if err != nil {
			return fmt.Errorf("socket Control: %w", err)
		}

		if cerr != nil {
			return fmt.Errorf("BindToDevice: %w", cerr)
		}

		return nil
	}
}

func switchNetwork(cur string) string {
	switch cur {
	case NetUDP:
		return NetTCP
	case NetUDP4:
		return NetTCP4
	case NetUDP6:
		return NetTCP6

	case NetTCP:
		return NetUDP
	case NetTCP4:
		return NetUDP4
	case NetTCP6:
		return NetUDP6

	default:
		return ""
	}
}
