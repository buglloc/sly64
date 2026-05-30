package upstream

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithConnPoolCanDisablePool(t *testing.T) {
	require.IsType(t, nopOpt{}, WithConnPool(0))
	require.IsType(t, nopOpt{}, WithConnPool(1))
	require.IsType(t, poolOpt{}, WithConnPool(2))
}

func TestParseAddrUsesJoinHostPortForIPv6(t *testing.T) {
	p, err := NewPlain(WithAddr("2001:db8::1", NetUDP))
	require.NoError(t, err)
	require.Equal(t, "udp://[2001:db8::1]:53", p.Address())

	d, err := NewDoT(WithAddr("2001:db8::1", NetTCPTLS))
	require.NoError(t, err)
	require.Equal(t, "tcp-tls://[2001:db8::1]:853", d.Address())
}

func TestDefaultDialTimeoutIsApplied(t *testing.T) {
	d := NewDialer()
	require.Equal(t, DefaultDialTimeout, d.Timeout)

	d = NewDialer(WithDialTimeout(time.Second))
	require.Equal(t, time.Second, d.Timeout)
}
