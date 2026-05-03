package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/config"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/database"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/logger"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to config file")
	flag.Parse()

	cfg, err := config.NewOrderServiceConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to init config", err)
		os.Exit(1)
	}

	slogger, err := logger.NewSLogger(cfg)
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to init logger", err)
		os.Exit(1)
	}

	pool, err := database.NewPool(cfg)
	if err != nil {
		slogger.Error("failed to init new postgres pool", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := pool.Close(); err != nil {
			slogger.Error("failed to close postgres pool", "error", err)
		}
	}()

	slogger.Debug("infrastructure added")
	os.Exit(0)
}
