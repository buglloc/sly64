package router

import (
	"fmt"
	"math"
	"time"

	"github.com/karlseguin/ccache/v3"
	"github.com/miekg/dns"
)

type CacheCfg struct {
	Size   int
	MinTTL uint32
	MaxTTL uint32
}

type cacheFetcher func() (*dns.Msg, error)

type Cache interface {
	Fetch(req *dns.Msg, fetcher cacheFetcher) (*dns.Msg, error)
}

func NewCache(cfg CacheCfg) Cache {
	if cfg.Size == 0 {
		return &nopCache{}
	}

	return &memCache{
		cfg: cfg,
		cache: ccache.New(
			ccache.Configure[*cacheEntry]().
				MaxSize(int64(cfg.Size)),
		),
	}
}

type nopCache struct {
}

func (c *nopCache) Fetch(_ *dns.Msg, fetch cacheFetcher) (*dns.Msg, error) {
	return fetch()
}

type memCache struct {
	cfg   CacheCfg
	cache *ccache.Cache[*cacheEntry]
}

type cacheEntry struct {
	msg     *dns.Msg
	created time.Time
}

func (c *memCache) Fetch(req *dns.Msg, fetch cacheFetcher) (*dns.Msg, error) {
	key := c.reqKey(req)
	if len(key) == 0 {
		return fetch()
	}

	item := c.cache.Get(key)
	if item != nil && !item.Expired() {
		return c.withAgedTTL(item), nil
	}

	rr, err := fetch()
	if err != nil {
		return nil, err
	}

	if rr.Rcode != dns.RcodeSuccess || rr.Truncated {
		return rr, nil
	}

	ttl := c.ttl(req, rr)
	if ttl == 0 {
		return rr, nil
	}

	c.cache.Set(key, &cacheEntry{
		msg:     rr.Copy(),
		created: time.Now(),
	}, ttl)

	return rr, nil
}

func (c *memCache) reqKey(req *dns.Msg) string {
	if len(req.Question) != 1 {
		return ""
	}

	q := req.Question[0]

	return fmt.Sprintf("%d:%d:%s", q.Qclass, q.Qtype, normalizeDomain(q.Name))
}

func (c *memCache) ttl(req, msg *dns.Msg) time.Duration {
	ttl := c.cfg.MaxTTL
	hasAnswer := false
	qtype := req.Question[0].Qtype

	for _, a := range msg.Answer {
		ttl = min(ttl, a.Header().Ttl)

		if a.Header().Rrtype == qtype || qtype == dns.TypeANY {
			hasAnswer = true
		}
	}

	if !hasAnswer {
		for _, rr := range msg.Ns {
			if soa, ok := rr.(*dns.SOA); ok {
				return time.Duration(min(ttl, soa.Hdr.Ttl, soa.Minttl)) * time.Second
			}
		}

		if len(msg.Answer) == 0 {
			return 0
		}

		// Without an SOA, a CNAME-only answer is not evidence of NODATA.
	}

	if c.cfg.MinTTL != 0 && ttl < c.cfg.MinTTL {
		ttl = c.cfg.MinTTL
	}

	return time.Duration(ttl) * time.Second
}

func (c *memCache) withAgedTTL(item *ccache.Item[*cacheEntry]) *dns.Msg {
	entry := item.Value()
	msg := entry.msg.Copy()
	ageSeconds := uint32(math.Ceil(time.Since(entry.created).Seconds()))

	ageRRsTTL(msg.Answer, ageSeconds)
	ageRRsTTL(msg.Ns, ageSeconds)
	ageRRsTTL(msg.Extra, ageSeconds)

	return msg
}

func ageRRsTTL(rrs []dns.RR, age uint32) {
	for _, rr := range rrs {
		if rr.Header().Rrtype == dns.TypeOPT {
			continue
		}

		rr.Header().Ttl = ageTTL(rr.Header().Ttl, age)
	}
}

func ageTTL(ttl uint32, age uint32) uint32 {
	if ttl <= age {
		return 0
	}

	return ttl - age
}
