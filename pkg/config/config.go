package config

import libconfig "github.com/ekristen/libnuke/pkg/config"

func New(opts libconfig.Options) (*Config, error) {
	cfg, err := libconfig.New(opts)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	if err := c.Load(opts.Path); err != nil {
		return nil, err
	}

	c.Config = *cfg
	return c, nil
}

type Config struct {
	libconfig.Config `yaml:",inline"`
}
