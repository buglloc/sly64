package router

import (
	"net"
	"testing"
	"testing/synctest"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestCacheKeyIncludesClassAndNormalizesName(t *testing.T) {
	c := NewCache(CacheCfg{Size: 10, MaxTTL: 60})
	fetches := 0

	reqIN := newCacheReq("Example.COM.", dns.ClassINET)

	rspIN, err := c.Fetch(reqIN, func() (*dns.Msg, error) {
		fetches++

		return newCacheRsp(reqIN, "192.0.2.1", 60), nil
	})
	require.NoError(t, err)
	require.Equal(t, "192.0.2.1", rspIN.Answer[0].(*dns.A).A.String())

	rspINLower, err := c.Fetch(newCacheReq("example.com.", dns.ClassINET), func() (*dns.Msg, error) {
		fetches++

		return newCacheRsp(reqIN, "192.0.2.2", 60), nil
	})
	require.NoError(t, err)
	require.Equal(t, "192.0.2.1", rspINLower.Answer[0].(*dns.A).A.String())

	reqCH := newCacheReq("example.com.", dns.ClassCHAOS)

	rspCH, err := c.Fetch(reqCH, func() (*dns.Msg, error) {
		fetches++

		return newCacheRsp(reqCH, "192.0.2.3", 60), nil
	})
	require.NoError(t, err)
	require.Equal(t, "192.0.2.3", rspCH.Answer[0].(*dns.A).A.String())
	require.Equal(t, 2, fetches)
}

func TestCacheAgesTTL(t *testing.T) {
	c := NewCache(CacheCfg{Size: 10, MaxTTL: 60})
	req := newCacheReq("example.com.", dns.ClassINET)

	_, err := c.Fetch(req, func() (*dns.Msg, error) {
		return newCacheRsp(req, "192.0.2.1", 60), nil
	})
	require.NoError(t, err)

	time.Sleep(1100 * time.Millisecond)

	rsp, err := c.Fetch(req, func() (*dns.Msg, error) {
		require.FailNow(t, "unexpected fetch")

		return nil, nil
	})
	require.NoError(t, err)
	require.Less(t, rsp.Answer[0].Header().Ttl, uint32(60))
}

func TestCacheNODATALifetime(t *testing.T) {
	tests := []struct {
		name      string
		soaTTL    uint32
		minimum   uint32
		minTTL    uint32
		maxTTL    uint32
		cnameTTLs []uint32
		signed    bool
		lifetime  time.Duration
	}{
		{
			name: "SOA TTL below MINIMUM", soaTTL: 2, minimum: 20, maxTTL: 60,
			lifetime: 2 * time.Second,
		},
		{
			name: "MINIMUM below SOA TTL", soaTTL: 20, minimum: 2, maxTTL: 60,
			lifetime: 2 * time.Second,
		},
		{
			name: "configured maximum", soaTTL: 20, minimum: 30, minTTL: 30, maxTTL: 2,
			lifetime: 2 * time.Second,
		},
		{
			name: "configured minimum does not extend NODATA", soaTTL: 20, minimum: 2, minTTL: 30, maxTTL: 60,
			lifetime: 2 * time.Second,
		},
		{
			name: "CNAME chain caps NODATA", soaTTL: 20, minimum: 30, minTTL: 10, maxTTL: 60,
			cnameTTLs: []uint32{10, 2}, lifetime: 2 * time.Second,
		},
		{
			name: "signed CNAME NODATA", soaTTL: 20, minimum: 2, minTTL: 30, maxTTL: 60,
			cnameTTLs: []uint32{20}, signed: true, lifetime: 2 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				c := NewCache(CacheCfg{Size: 10, MinTTL: tt.minTTL, MaxTTL: tt.maxTTL}).(*memCache)

				defer c.cache.Stop()

				req := newCacheReq("example.com.", dns.ClassINET)
				fetches := 0

				fetch := func() (*dns.Msg, error) {
					fetches++

					rsp := newCacheNODATARsp(req, tt.soaTTL, tt.minimum)
					rsp.Ns[0].(*dns.SOA).Serial = uint32(fetches)
					rsp.Answer = newCacheCNAMEChain(req, tt.cnameTTLs)

					if tt.signed {
						rsp.Answer = append(rsp.Answer, &dns.RRSIG{
							Hdr:         dns.RR_Header{Name: req.Question[0].Name, Rrtype: dns.TypeRRSIG, Class: dns.ClassINET, Ttl: 20},
							TypeCovered: dns.TypeCNAME,
						})
					}

					return rsp, nil
				}

				rsp, err := c.Fetch(req, fetch)
				require.NoError(t, err)
				require.Equal(t, uint32(1), rsp.Ns[0].(*dns.SOA).Serial)

				time.Sleep(tt.lifetime - time.Nanosecond)

				rsp, err = c.Fetch(req, fetch)
				require.NoError(t, err)
				require.Equal(t, uint32(1), rsp.Ns[0].(*dns.SOA).Serial, "NODATA should still be cached before expiry")
				require.Equal(t, 1, fetches)

				time.Sleep(2 * time.Nanosecond)

				rsp, err = c.Fetch(req, fetch)
				require.NoError(t, err)
				require.Equal(t, uint32(2), rsp.Ns[0].(*dns.SOA).Serial, "expired NODATA must be refetched")
				require.Equal(t, 2, fetches)
			})
		})
	}
}

func TestCacheDoesNotCacheUncacheableResponses(t *testing.T) {
	tests := []struct {
		name      string
		soa       bool
		soaTTL    uint32
		minimum   uint32
		minTTL    uint32
		maxTTL    uint32
		cnameTTLs []uint32
		positive  bool
		rcode     int
	}{
		{name: "empty answer without SOA", minTTL: 30, maxTTL: 60},
		{name: "zero SOA TTL", soa: true, minimum: 20, minTTL: 30, maxTTL: 60},
		{name: "zero SOA MINIMUM", soa: true, soaTTL: 20, minTTL: 30, maxTTL: 60},
		{name: "zero CNAME TTL in NODATA", soa: true, soaTTL: 20, minimum: 20, minTTL: 30, maxTTL: 60, cnameTTLs: []uint32{0}},
		{name: "zero positive TTL", positive: true, maxTTL: 60},
		{name: "zero configured maximum", soa: true, soaTTL: 20, minimum: 20, minTTL: 30},
		{name: "NXDOMAIN remains uncached", soa: true, soaTTL: 20, minimum: 20, maxTTL: 60, rcode: dns.RcodeNameError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				c := NewCache(CacheCfg{Size: 10, MinTTL: tt.minTTL, MaxTTL: tt.maxTTL}).(*memCache)

				defer c.cache.Stop()

				req := newCacheReq("example.com.", dns.ClassINET)
				fetches := 0

				fetch := func() (*dns.Msg, error) {
					fetches++

					rsp := new(dns.Msg)
					rsp.SetReply(req)

					if tt.soa {
						rsp = newCacheNODATARsp(req, tt.soaTTL, tt.minimum)
					}

					rsp.Answer = newCacheCNAMEChain(req, tt.cnameTTLs)

					if tt.positive {
						rsp = newCacheRsp(req, "192.0.2.1", 0)
					}

					rsp.Rcode = tt.rcode

					return rsp, nil
				}

				for range 2 {
					rsp, err := c.Fetch(req, fetch)
					require.NoError(t, err)
					require.Equal(t, tt.rcode, rsp.Rcode)
				}

				require.Equal(t, 2, fetches, "each request must reach the upstream")
			})
		})
	}
}

func TestCachePositiveMinTTL(t *testing.T) {
	tests := []struct {
		name  string
		qtype uint16
		cname bool
		soa   bool
	}{
		{name: "address answer", qtype: dns.TypeA},
		{name: "address answer with SOA", qtype: dns.TypeA, soa: true},
		{name: "CNAME-only answer without SOA", qtype: dns.TypeA, cname: true},
		{name: "requested CNAME with SOA", qtype: dns.TypeCNAME, cname: true, soa: true},
		{name: "ANY answer with CNAME and SOA", qtype: dns.TypeANY, cname: true, soa: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				c := NewCache(CacheCfg{Size: 10, MinTTL: 5, MaxTTL: 60}).(*memCache)

				defer c.cache.Stop()

				req := newCacheReq("example.com.", dns.ClassINET)
				req.Question[0].Qtype = tt.qtype
				fetches := 0

				fetch := func() (*dns.Msg, error) {
					fetches++

					rsp := newCacheRsp(req, "192.0.2.1", 1)

					if tt.cname {
						rsp.Answer = newCacheCNAMEChain(req, []uint32{1})
					}

					if tt.soa {
						rsp.Ns = newCacheNODATARsp(req, 1, 1).Ns
					}

					return rsp, nil
				}

				_, err := c.Fetch(req, fetch)
				require.NoError(t, err)

				time.Sleep(5*time.Second - time.Nanosecond)

				rsp, err := c.Fetch(req, fetch)
				require.NoError(t, err)
				require.Equal(t, 1, fetches, "positive MinTTL must still extend cache lifetime")
				require.Equal(t, uint32(0), rsp.Answer[0].Header().Ttl)

				time.Sleep(2 * time.Nanosecond)

				rsp, err = c.Fetch(req, fetch)
				require.NoError(t, err)
				require.Equal(t, 2, fetches)
				require.Equal(t, uint32(1), rsp.Answer[0].Header().Ttl)
			})
		})
	}
}

func newCacheNODATARsp(req *dns.Msg, ttl, minimum uint32) *dns.Msg {
	rsp := new(dns.Msg)
	rsp.SetReply(req)

	rsp.Ns = []dns.RR{&dns.SOA{
		Hdr:    dns.RR_Header{Name: "example.com.", Rrtype: dns.TypeSOA, Class: req.Question[0].Qclass, Ttl: ttl},
		Ns:     "ns.example.com.",
		Mbox:   "hostmaster.example.com.",
		Minttl: minimum,
	}}

	return rsp
}

func newCacheCNAMEChain(req *dns.Msg, ttls []uint32) []dns.RR {
	var answers []dns.RR
	name := req.Question[0].Name

	for i, ttl := range ttls {
		target := "alias.example.com."
		if i > 0 {
			target = "terminal.example.com."
		}

		answers = append(answers, &dns.CNAME{
			Hdr:    dns.RR_Header{Name: name, Rrtype: dns.TypeCNAME, Class: req.Question[0].Qclass, Ttl: ttl},
			Target: target,
		})

		name = target
	}

	return answers
}

func newCacheReq(name string, class uint16) *dns.Msg {
	return &dns.Msg{Question: []dns.Question{{Name: name, Qtype: dns.TypeA, Qclass: class}}}
}

func newCacheRsp(req *dns.Msg, ip string, ttl uint32) *dns.Msg {
	rsp := new(dns.Msg)
	rsp.SetReply(req)

	rsp.Answer = []dns.RR{&dns.A{
		Hdr: dns.RR_Header{Name: req.Question[0].Name, Rrtype: dns.TypeA, Class: req.Question[0].Qclass, Ttl: ttl},
		A:   net.ParseIP(ip),
	}}

	return rsp
}
