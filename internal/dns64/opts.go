package dns64

import (
	"fmt"
	"net"
)

type Option func(*DNS64) error

func WithPrefix(addr string) Option {
	return func(d *DNS64) error {
		_, pref, err := net.ParseCIDR(addr)
		if err != nil {
			return fmt.Errorf("invalid prefix: %w", err)
		}

		// Test for valid prefix
		ones, bits := pref.Mask.Size()
		if bits != 128 {
			return fmt.Errorf("invalid netmask %d IPv6 address: %q", bits, pref)
		}

		if ones%8 != 0 || ones < 32 || ones > 96 {
			return fmt.Errorf("invalid prefix length (required >=32 && <96): %q", pref)
		}

		d.prefix = pref
		return nil
	}
}
