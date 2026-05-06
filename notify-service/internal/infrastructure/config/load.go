package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func NewNotifyServiceConfig(configPath string) (*OMSNotifyServiceConfig, error) {
	cfg := new(OMSNotifyServiceConfig)
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read environments: %w", err)
	}

	return cfg, nil
}
