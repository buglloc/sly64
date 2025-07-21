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

func NewRuntime(configs ...string) (*Runtime, error) {
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

	if err := loadConfigs(cfg, configs...); err != nil {
		return nil, err
	}

	return newRuntime(cfg)
}

func loadConfigs(cfg *configpb.Config, paths ...string) error {
	for _, p := range paths {
		if p == "" {
			continue
		}

		if err := loadConfig(cfg, p); err != nil {
			return fmt.Errorf("load config %q: %w", p, err)
		}
	}

	return nil
}

func loadConfig(cfg *configpb.Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	return prototext.Unmarshal(data, cfg)
}
