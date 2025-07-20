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

	netUDP  = "udp"
	netUDP4 = "udp4"
	netUDP6 = "udp6"

	netTCP  = "tcp"
	netTCP4 = "tcp4"
	netTCP6 = "tcp6"
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
	case netUDP:
		return netTCP
	case netUDP4:
		return netTCP4
	case netUDP6:
		return netTCP6

	case netTCP:
		return netUDP
	case netTCP4:
		return netUDP4
	case netTCP6:
		return netUDP6

	default:
		return ""
	}
}
