package config

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"

	"github.com/buglloc/sly64/v2/internal/upstream"
)

var ShutdownDeadline = 1 * time.Minute

type Config struct {
	LogLevel    zerolog.Level `yaml:"log_level"`
	MaxRequests int           `yaml:"max_requests"`
	Listeners   []Listener    `yaml:"listeners"`
	Routes      []Route       `yaml:"routes"`
}

func (c *Config) NewRuntime() (*Runtime, error) {
	return newRuntime(c)
}

func (c *Config) Validate() error {
	return nil
}

func Load(configs ...string) (*Config, error) {
	cfg := &Config{
		LogLevel: zerolog.InfoLevel,
		Listeners: []Listener{
			{
				Addr: ":1353",
				Net:  "udp",
			},
			{
				Addr: ":1353",
				Net:  "tcp",
			},
		},
		Routes: []Route{
			{
				Name: "default",
				Upstreams: []Upstream{
					{
						Kind: upstream.KindPlain,
						Plain: PlainUpstream{
							Addr: "udp://1.1.1.1:53",
						},
					},
					{
						Kind: upstream.KindPlain,
						Plain: PlainUpstream{
							Addr: "udp://1.0.0.1:53",
						},
					},
				},
				Domains: []string{
					"*.",
				},
			},
		},
	}

	if err := loadConfigs(cfg, configs...); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func loadConfigs(cfg *Config, paths ...string) error {
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

func loadConfig(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	return yaml.UnmarshalWithOptions(data, cfg, yaml.DisallowUnknownField())
}
