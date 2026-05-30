package router

import (
	"net"
	"testing"
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
