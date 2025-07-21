package dns64_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buglloc/sly64/v2/internal/dns64"
)

func TestTo6(t *testing.T) {
	cases := []struct {
		prefix   string
		addr     string
		expected string
	}{
		{
			prefix:   "fd00:cafe:babe:64:1::/96",
			addr:     "64.64.64.64",
			expected: "fd00:cafe:babe:64:1:0:4040:4040",
		},
		{
			prefix:   "64:ff9b::/96",
			addr:     "64.64.64.64",
			expected: "64:ff9b::4040:4040",
		},
		{
			prefix:   "64:ff9b::/64",
			addr:     "64.64.64.64",
			expected: "64:ff9b::40:4040:4000:0",
		},
		{
			prefix:   "64:ff9b::/56",
			addr:     "64.64.64.64",
			expected: "64:ff9b:0:40:40:4040::",
		},
		{
			prefix:   "64::/32",
			addr:     "64.64.64.64",
			expected: "64:0:4040:4040::",
		},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s@%s", c.addr, c.prefix)
		t.Run(name, func(t *testing.T) {
			d64, err := dns64.NewDNS64(dns64.WithPrefix(c.prefix))
			require.NoError(t, err)

			ip, err := d64.To6(net.ParseIP(c.addr))
			require.NoError(t, err)
			require.Equal(t, c.expected, ip.String())
		})
	}
}
