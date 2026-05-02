package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/identicalaffiliation/oms-with-events/migrator/internal/config"
	"github.com/identicalaffiliation/oms-with-events/migrator/internal/logger"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to config file")
	flag.Parse()

	config, err := config.NewMigratorConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "load config: %v", err)
		os.Exit(1)
	}

	logger, err := logger.NewSLogger(config)
	if err != nil {
		fmt.Fprintf(os.Stdout, "init logger: %v", err)
		os.Exit(1)
	}

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.PostgresConfig.Username,
		config.PostgresConfig.Password,
		config.PostgresConfig.Host,
		config.PostgresConfig.Port,
		config.PostgresConfig.DBName,
		config.PostgresConfig.SSLMode)

	m, err := migrate.New(config.MigrationsPath, dataSourceName)
	if err != nil {
		logger.Error("new migrate failed", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Debug("no new migrations")
			logger.Debug("success")
			os.Exit(0)
		}

		logger.Error("up migrations failed", "error", err)
		os.Exit(1)
	}

	logger.Debug("migrations add successfuly")
	os.Exit(0)
}
