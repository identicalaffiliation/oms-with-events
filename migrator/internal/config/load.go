package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func NewMigratorConfig(configPath string) (*OMSMigratorConfig, error) {
	cfg := new(OMSMigratorConfig)
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read environments: %w", err)
	}

	return cfg, nil
}
