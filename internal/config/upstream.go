package config

import (
	"fmt"
	"time"

	"github.com/buglloc/sly64/v2/internal/upstream"
)

type Upstream struct {
	Kind  upstream.Kind `yaml:"kind"`
	Plain PlainUpstream `yaml:"plain"`
}

type PlainUpstream struct {
	Addr        string        `yaml:"addr"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
	Interface   string        `yaml:"iface"`
}

func NewUpstream(cfg Upstream) (upstream.Upstream, error) {
	switch cfg.Kind {
	case upstream.KindPlain:
		return upstream.NewPlain(
			upstream.WithAddr(cfg.Plain.Addr),
			upstream.WithDialTimeout(cfg.Plain.DialTimeout),
			upstream.WithTimeout(cfg.Plain.Timeout),
			upstream.WithIface(cfg.Plain.Interface),
		)

	default:
		return nil, fmt.Errorf("unsupported upstream: %s", cfg.Kind)
	}
}
