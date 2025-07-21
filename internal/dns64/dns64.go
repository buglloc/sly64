package dns64

import (
	"fmt"
	"net"
)

const (
	DefaultPrefix = "64:ff9b::/96"
)

type DNS64 struct {
	prefix *net.IPNet
}

func NewDNS64(opts ...Option) (*DNS64, error) {
	n := &DNS64{}
	_ = WithPrefix(DefaultPrefix)(n)

	for _, opt := range opts {
		if err := opt(n); err != nil {
			return nil, fmt.Errorf("invalid option: %w", err)
		}
	}

	return n, nil
}

// To6 takes a prefix and IPv4 address and returns an IPv6 address according to RFC 6052.
func (d *DNS64) To6(addr net.IP) (net.IP, error) {
	addr = addr.To4()
	if addr == nil {
		return nil, fmt.Errorf("not a valid IPv4 address: %s", addr)
	}

	n, _ := d.prefix.Mask.Size()
	// Assumes prefix has been validated during setup
	v6 := make([]byte, 16)
	i, j := 0, 0

	for ; i < n/8; i++ {
		v6[i] = d.prefix.IP[i]
	}
	for ; i < 8; i, j = i+1, j+1 {
		v6[i] = addr[j]
	}
	if i == 8 {
		i++
	}
	for ; j < 4; i, j = i+1, j+1 {
		v6[i] = addr[j]
	}

	return v6, nil
}
