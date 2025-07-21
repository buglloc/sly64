package router

import (
	"fmt"
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
			ccache.Configure[*dns.Msg]().
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
	cache *ccache.Cache[*dns.Msg]
}

func (c *memCache) Fetch(req *dns.Msg, fetch cacheFetcher) (*dns.Msg, error) {
	key := c.reqKey(req)
	if len(key) == 0 {
		return fetch()
	}

	item := c.cache.Get(key)
	if item != nil && !item.Expired() {
		return item.Value().Copy(), nil
	}

	rr, err := fetch()
	if err != nil {
		return nil, err
	}

	if rr.Rcode != dns.RcodeSuccess || rr.Truncated {
		return rr, nil
	}

	c.cache.Set(key, rr.Copy(), c.ttl(rr))
	return rr, nil
}

func (c *memCache) reqKey(req *dns.Msg) string {
	if len(req.Question) != 1 {
		return ""
	}

	return fmt.Sprintf("%d:%s", req.Question[0].Qtype, req.Question[0].Name)
}

func (c *memCache) ttl(msg *dns.Msg) time.Duration {
	ttl := c.cfg.MaxTTL
	for _, a := range msg.Answer {
		aTTL := a.Header().Ttl
		if aTTL < ttl {
			ttl = aTTL
		}
	}

	if c.cfg.MinTTL != 0 && ttl < c.cfg.MinTTL {
		ttl = c.cfg.MinTTL
	}

	return time.Duration(ttl) * time.Second
}
