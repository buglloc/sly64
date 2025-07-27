package config

import (
	"fmt"
	"os"
	"time"

	"google.golang.org/protobuf/encoding/prototext"

	"github.com/buglloc/sly64/v2/internal/config/configpb"
)

const (
	ShutdownDeadline = 1 * time.Minute
)

type cfgPatcher func(cfg *configpb.Config, path string) error

var cfgPatchers = []cfgPatcher{
	routerPatcher,
	upstreamPatcher,
}

func NewRuntime(cfgPath string) (*Runtime, error) {
	cfg := &configpb.Config{
		LogLevel: configpb.LogLevel_LOG_LEVEL_INFO,
		Listener: []*configpb.Listener{
			{
				Addr: ":1353",
				Net:  configpb.Net_NET_UDP,
			},
			{
				Addr: ":1353",
				Net:  configpb.Net_NET_TCP,
			},
		},
		Route: []*configpb.Route{
			{
				Name: "default",
				Upstream: []*configpb.Upstream{
					{
						Kind: &configpb.Upstream_Udp{
							Udp: &configpb.UdpUpstream{
								Addr: "1.1.1.1:53",
							},
						},
					},
					{
						Kind: &configpb.Upstream_Tcp{
							Tcp: &configpb.TcpUpstream{
								Addr: "1.1.1.1:53",
							},
						},
					},
				},
				Source: []*configpb.Source{
					{
						Kind: &configpb.Source_Static{
							Static: &configpb.StaticSource{
								Domain: []string{
									"*.",
								},
							},
						},
					},
				},
			},
		},
	}

	if len(cfgPath) > 0 {
		if err := loadConfig(cfg, cfgPath); err != nil {
			return nil, fmt.Errorf("load config %q: %w", cfgPath, err)
		}
	}

	return newRuntime(cfg)
}

func loadConfig(cfg *configpb.Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	if err := prototext.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	for _, p := range cfgPatchers {
		if err := p(cfg, path); err != nil {
			return fmt.Errorf("patch config: %w", err)
		}
	}

	return nil
}
