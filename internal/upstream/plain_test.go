package upstream

import (
	"context"
	"encoding/binary"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestPlainQueryTimeoutIndependentOfDialTimeout(t *testing.T) {
	for _, tc := range []struct {
		name     string
		fallback bool
	}{
		{name: "UDP"},
		{name: "fallback TCP", fallback: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const dialTimeout = 50 * time.Millisecond
			received := make(chan struct{}, 1)
			release := make(chan struct{})
			addr := startPlainTestServer(t, func(ctx context.Context, network string, req *dns.Msg) ([]byte, error) {
				if tc.fallback && network == NetUDP {
					rsp := new(dns.Msg)
					rsp.SetReply(req)
					rsp.Truncated = true

					return rsp.Pack()
				}

				received <- struct{}{}
				select {
				case <-release:
					return plainTestReply(req).Pack()
				case <-ctx.Done():
					return nil, nil
				}
			})

			p, err := NewPlain(
				WithAddr(addr, NetUDP4),
				WithDialTimeout(dialTimeout),
				WithTimeout(time.Second),
			)
			require.NoError(t, err)

			t.Cleanup(func() { require.NoError(t, p.Close()) })

			req := plainTestQuery()
			ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)

			type result struct {
				rsp *dns.Msg
				err error
			}

			results := make(chan result, 1)
			var client sync.WaitGroup

			client.Go(func() {
				rsp, err := p.Exchange(ctx, req)
				results <- result{rsp: rsp, err: err}
			})

			t.Cleanup(func() {
				cancel()
				client.Wait()
			})

			select {
			case <-received:
			case <-ctx.Done():
				t.Fatal("responder did not receive the query")
			}

			// Release the answer only after the dial timeout has elapsed on
			// an already-established socket, but before the query timeout.
			delay := time.NewTimer(3 * dialTimeout)
			defer delay.Stop()

			select {
			case got := <-results:
				t.Fatalf("query ended before the answer was released: %v", got.err)
			case <-delay.C:
			}

			close(release)

			select {
			case got := <-results:
				require.NoError(t, got.err)
				require.Equal(t, req.Id, got.rsp.Id)
				requireDNSAnswers(t, plainTestReply(req).Answer, got.rsp.Answer)
			case <-ctx.Done():
				t.Fatal("query did not receive the delayed answer")
			}
		})
	}
}

func TestPlainQueryTimeoutStillEnforced(t *testing.T) {
	for _, tc := range []struct {
		name     string
		fallback bool
	}{
		{name: "UDP"},
		{name: "fallback TCP", fallback: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			received := make(chan struct{}, 1)
			addr := startPlainTestServer(t, func(ctx context.Context, network string, req *dns.Msg) ([]byte, error) {
				if tc.fallback && network == NetUDP {
					rsp := new(dns.Msg)
					rsp.SetReply(req)
					rsp.Truncated = true

					return rsp.Pack()
				}

				received <- struct{}{}
				<-ctx.Done()

				return nil, nil
			})

			p, err := NewPlain(
				WithAddr(addr, NetUDP4),
				WithDialTimeout(time.Second),
				WithTimeout(100*time.Millisecond),
			)
			require.NoError(t, err)

			t.Cleanup(func() { require.NoError(t, p.Close()) })

			ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
			defer cancel()

			started := time.Now()
			rsp, err := p.Exchange(ctx, plainTestQuery())

			require.Nil(t, rsp)
			var netErr net.Error
			require.ErrorAs(t, err, &netErr)
			require.True(t, netErr.Timeout(), "query error is not a timeout: %v", err)
			require.NoError(t, ctx.Err(), "caller deadline, rather than query timeout, ended the exchange")
			require.Less(t, time.Since(started), time.Second)

			select {
			case <-received:
			default:
				t.Fatal("responder did not receive the query")
			}
		})
	}
}

func TestPlainMalformedResponseFallback(t *testing.T) {
	for _, tc := range []struct {
		name            string
		truncated       bool
		disableFallback bool
		mismatchedID    bool
		wantFallback    bool
	}{
		{name: "truncated decode error", truncated: true, wantFallback: true},
		{name: "truncated decode error with fallback disabled", truncated: true, disableFallback: true},
		{name: "non-truncated decode error"},
		{name: "truncated decode error with mismatched ID", truncated: true, mismatchedID: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			addr := startPlainTestServer(t, func(_ context.Context, network string, req *dns.Msg) ([]byte, error) {
				if network == NetTCP {
					return plainTestReply(req).Pack()
				}

				rsp := new(dns.Msg)
				rsp.SetReply(req)
				rsp.Truncated = tc.truncated
				if tc.mismatchedID {
					rsp.Id ^= 1
				}

				packet, err := rsp.Pack()
				if err != nil {
					return nil, err
				}

				// Advertise one answer but end after its compressed owner name,
				// preserving the question and TC bit in the partial response.
				binary.BigEndian.PutUint16(packet[6:8], 1)

				return append(packet, 0xc0, 0x0c), nil
			})

			p, err := NewPlain(
				WithAddr(addr, NetUDP4),
				WithTCPFallback(!tc.disableFallback),
				WithTimeout(time.Second),
			)
			require.NoError(t, err)

			t.Cleanup(func() { require.NoError(t, p.Close()) })

			ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
			defer cancel()

			req := plainTestQuery()

			rsp, err := p.Exchange(ctx, req)

			if tc.wantFallback {
				require.NoError(t, err)
				require.Equal(t, req.Id, rsp.Id)
				require.False(t, rsp.Truncated)
				requireDNSAnswers(t, plainTestReply(req).Answer, rsp.Answer)

				return
			}

			require.Error(t, err)
			require.NotNil(t, rsp, "decode error discarded its partial response")
			require.Equal(t, tc.truncated, rsp.Truncated)
			require.Equal(t, req.Question, rsp.Question)

			wantID := req.Id
			if tc.mismatchedID {
				wantID ^= 1
			}

			require.Equal(t, wantID, rsp.Id)
		})
	}
}

func TestPlainTruncatedAuthenticationErrorDoesNotFallback(t *testing.T) {
	addr := startPlainTestServer(t, func(_ context.Context, network string, req *dns.Msg) ([]byte, error) {
		if network == NetTCP {
			return plainTestReply(req).Pack()
		}

		rsp := new(dns.Msg)
		rsp.SetReply(req)
		rsp.Truncated = true
		rsp.SetTsig("key.test.", dns.HmacSHA256, 300, time.Now().Unix())

		return rsp.Pack()
	})

	p, err := NewPlain(WithAddr(addr, NetUDP4), WithTimeout(time.Second))
	require.NoError(t, err)

	t.Cleanup(func() { require.NoError(t, p.Close()) })

	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()

	req := plainTestQuery()

	rsp, err := p.Exchange(ctx, req)

	require.ErrorIs(t, err, dns.ErrSecret)
	require.NotNil(t, rsp)
	require.Equal(t, req.Id, rsp.Id)
	require.True(t, rsp.Truncated)
	require.NotNil(t, rsp.IsTsig())
}

func TestPlainUDPReceiveBufferFromEDNS(t *testing.T) {
	answer := []dns.RR{&dns.TXT{
		Hdr: dns.RR_Header{Name: "example.test.", Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 60},
		Txt: []string{strings.Repeat("a", 200), strings.Repeat("b", 200), strings.Repeat("c", 200)},
	}}

	addr := startPlainTestServer(t, func(_ context.Context, _ string, req *dns.Msg) ([]byte, error) {
		rsp := new(dns.Msg)
		rsp.SetReply(req)
		rsp.Answer = answer

		return rsp.Pack()
	})

	p, err := NewPlain(WithAddr(addr, NetUDP4), WithTCPFallback(false), WithTimeout(time.Second))
	require.NoError(t, err)

	t.Cleanup(func() { require.NoError(t, p.Close()) })

	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()

	req := new(dns.Msg)
	req.SetQuestion("example.test.", dns.TypeTXT)
	req.SetEdns0(4096, false)

	rsp, err := p.Exchange(ctx, req)

	require.NoError(t, err)
	require.Equal(t, req.Id, rsp.Id)
	requireDNSAnswers(t, answer, rsp.Answer)
}

func startPlainTestServer(t *testing.T, respond func(context.Context, string, *dns.Msg) ([]byte, error)) string {
	t.Helper()

	packetConn, err := net.ListenPacket(NetUDP4, "127.0.0.1:0")
	require.NoError(t, err)

	listener, err := net.Listen(NetTCP4, packetConn.LocalAddr().String())
	if err != nil {
		_ = packetConn.Close()
		t.Fatalf("listen TCP: %v", err)
	}

	started := make(chan struct{}, 2)
	finished := make(chan error, 2)
	responseErrors := make(chan error, 8)

	makeHandler := func(network string) dns.Handler {
		return dns.HandlerFunc(func(w dns.ResponseWriter, req *dns.Msg) {
			packet, err := respond(t.Context(), network, req)
			if err == nil && packet != nil {
				_, err = w.Write(packet)
			}

			if err != nil {
				responseErrors <- err
			}
		})
	}

	servers := []*dns.Server{
		{Net: NetUDP, PacketConn: packetConn, Handler: makeHandler(NetUDP)},
		{Net: NetTCP, Listener: listener, Handler: makeHandler(NetTCP)},
	}

	for _, server := range servers {
		server.NotifyStartedFunc = func() { started <- struct{}{} }
		go func() { finished <- server.ActivateAndServe() }()
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		for _, server := range servers {
			require.NoError(t, server.ShutdownContext(ctx))
		}

		for range servers {
			require.NoError(t, <-finished)
		}

		close(responseErrors)

		for err := range responseErrors {
			require.NoError(t, err)
		}
	})

	for range servers {
		select {
		case <-started:
		case err := <-finished:
			t.Fatalf("DNS responder stopped before startup: %v", err)
		case <-time.After(3 * time.Second):
			t.Fatal("DNS responder did not start")
		}
	}

	return packetConn.LocalAddr().String()
}

func plainTestQuery() *dns.Msg {
	req := new(dns.Msg)
	req.SetQuestion("example.test.", dns.TypeA)

	return req
}

func plainTestReply(req *dns.Msg) *dns.Msg {
	rsp := new(dns.Msg)
	rsp.SetReply(req)
	rsp.Answer = []dns.RR{&dns.A{
		Hdr: dns.RR_Header{Name: req.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
		A:   net.IP{192, 0, 2, 2},
	}}

	return rsp
}

func requireDNSAnswers(t *testing.T, want, got []dns.RR) {
	t.Helper()

	require.Len(t, got, len(want))

	for i := range want {
		require.Equal(t, want[i].String(), got[i].String())
	}
}
