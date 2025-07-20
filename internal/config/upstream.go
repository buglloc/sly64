package config

import (
	"fmt"
	"time"

	"github.com/buglloc/sly64/v2/internal/upstream"
)

type Upstream struct {
	Kind  upstream.Kind `yaml:"kind"`
	Plain PlainUpstream `yaml:"plain"`
	DoT   DoTUpstream   `yaml:"dot"`
}

type PlainUpstream struct {
	Addr        string        `yaml:"addr"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
	Interface   string        `yaml:"iface"`
}

type DoTUpstream struct {
	Addr        string        `yaml:"addr"`
	ServerName  string        `yaml:"server_name"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
	Interface   string        `yaml:"iface"`
}

func NewUpstream(cfg Upstream) (upstream.Upstream, error) {
	switch cfg.Kind {
	case upstream.KindPlain:
		return upstream.NewPlain(
			upstream.WithPlainAddr(cfg.Plain.Addr),
			upstream.WithDialTimeout(cfg.Plain.DialTimeout),
			upstream.WithTimeout(cfg.Plain.Timeout),
			upstream.WithIface(cfg.Plain.Interface),
		)

	case upstream.KindDoT:
		return upstream.NewDoT(
			upstream.WithDoTAddr(cfg.DoT.Addr, cfg.DoT.ServerName),
			upstream.WithDialTimeout(cfg.DoT.DialTimeout),
			upstream.WithTimeout(cfg.DoT.Timeout),
			upstream.WithIface(cfg.DoT.Interface),
		)

	default:
		return nil, fmt.Errorf("unsupported upstream: %s", cfg.Kind)
	}
}
