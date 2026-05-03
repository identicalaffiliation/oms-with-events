package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func NewOrderServiceConfig(configPath string) (*OMSGOrderServiceConfig, error) {
	cfg := new(OMSGOrderServiceConfig)
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read environments: %w", err)
	}

	return cfg, nil
}
