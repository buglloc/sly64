package listener

import (
	"net"
	"strings"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

type fakeWriter struct {
	remote net.Addr
	msg    *dns.Msg
}

func (f *fakeWriter) WriteMsg(m *dns.Msg) error   { f.msg = m; return nil }
func (f *fakeWriter) Write(b []byte) (int, error) { return len(b), nil }
func (f *fakeWriter) Close() error                { return nil }
func (f *fakeWriter) TsigStatus() error           { return nil }
func (f *fakeWriter) TsigTimersOnly(bool)         {}
func (f *fakeWriter) Hijack()                     {}
func (f *fakeWriter) LocalAddr() net.Addr         { return nil }
func (f *fakeWriter) RemoteAddr() net.Addr        { return f.remote }
func (f *fakeWriter) Network() string             { return f.remote.Network() }

func bigTXTAnswer(count int) []dns.RR {
	out := make([]dns.RR, 0, count)
	for range count {
		out = append(out, &dns.TXT{
			Hdr: dns.RR_Header{
				Name:   "example.com.",
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			Txt: []string{strings.Repeat("a", 100)},
		})
	}
	return out
}

func countOPT(rrs []dns.RR) int {
	n := 0
	for _, rr := range rrs {
		if rr.Header().Rrtype == dns.TypeOPT {
			n++
		}
	}
	return n
}

func TestServer_writeMsg(t *testing.T) {
	udpAddr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 5353}
	tcpAddr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 5353}

	type clientEDNS0 struct {
		udpSize uint16
		do      bool
	}

	cases := []struct {
		name             string
		transport        net.Addr
		clientEDNS0      *clientEDNS0
		answers          int
		upstreamOPTSize  uint16 // 0 means: no upstream OPT in rsp.Extra
		wantTC           bool
		wantOPTCount     int
		wantOPTUDPSize   uint16
		wantOPTDo        bool
		wantAnswersExact int // -1 means: skip the check (e.g. truncated UDP)
		wantPackedAtMost int // 0 means: skip the check
	}{
		{
			name:             "udp/no edns0/large response truncated to 512",
			transport:        udpAddr,
			answers:          10,
			wantTC:           true,
			wantOPTCount:     0,
			wantAnswersExact: -1,
			wantPackedAtMost: dns.MinMsgSize,
		},
		{
			name:             "udp/edns0 4096/large response fits without truncation",
			transport:        udpAddr,
			clientEDNS0:      &clientEDNS0{udpSize: 4096},
			answers:          10,
			wantOPTCount:     1,
			wantOPTUDPSize:   4096,
			wantAnswersExact: 10,
		},
		{
			name:             "udp/no edns0/upstream opt stripped",
			transport:        udpAddr,
			upstreamOPTSize:  1232,
			wantOPTCount:     0,
			wantAnswersExact: 0,
		},
		{
			name:             "udp/edns0 1232+DO/upstream opt replaced by client values",
			transport:        udpAddr,
			clientEDNS0:      &clientEDNS0{udpSize: 1232, do: true},
			upstreamOPTSize:  4096,
			wantOPTCount:     1,
			wantOPTUDPSize:   1232,
			wantOPTDo:        true,
			wantAnswersExact: 0,
		},
		{
			name:             "tcp/no edns0/large response not truncated",
			transport:        tcpAddr,
			answers:          10,
			wantOPTCount:     0,
			wantAnswersExact: 10,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := new(dns.Msg)
			req.SetQuestion("example.com.", dns.TypeTXT)
			if c.clientEDNS0 != nil {
				req.SetEdns0(c.clientEDNS0.udpSize, c.clientEDNS0.do)
			}

			rsp := new(dns.Msg)
			rsp.SetReply(req)
			rsp.Answer = bigTXTAnswer(c.answers)
			if c.upstreamOPTSize > 0 {
				opt := &dns.OPT{
					Hdr: dns.RR_Header{Name: ".", Rrtype: dns.TypeOPT},
				}
				opt.SetUDPSize(c.upstreamOPTSize)
				rsp.Extra = append(rsp.Extra, opt)
			}

			w := &fakeWriter{remote: c.transport}
			(&Server{}).writeMsg(w, req, rsp)
			require.NotNil(t, w.msg)

			require.Equal(t, c.wantTC, w.msg.Truncated, "TC bit")
			require.Equal(t, c.wantOPTCount, countOPT(w.msg.Extra), "OPT count")

			if c.wantOPTCount == 1 {
				opt := w.msg.IsEdns0()
				require.NotNil(t, opt)
				require.Equal(t, c.wantOPTUDPSize, opt.UDPSize(), "OPT UDPSize")
				require.Equal(t, c.wantOPTDo, opt.Do(), "OPT DO bit")
			}

			if c.wantAnswersExact >= 0 {
				require.Len(t, w.msg.Answer, c.wantAnswersExact, "answer count")
			}

			if c.wantPackedAtMost > 0 {
				data, err := w.msg.Pack()
				require.NoError(t, err)
				require.LessOrEqual(t, len(data), c.wantPackedAtMost, "packed length")
			}
		})
	}
}
