package database

import (
	"database/sql"
	"fmt"

	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/config"
	_ "github.com/lib/pq"
)

const (
	POSTGRES_DRIVER_NAME = "postgres"
)

func NewPool(cfg *config.OMSGOrderServiceConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresConfig.Host, cfg.PostgresConfig.Port, cfg.PostgresConfig.Username,
		cfg.PostgresConfig.Password, cfg.PostgresConfig.DBName, cfg.PostgresConfig.SSLMode)

	pool, err := sql.Open(POSTGRES_DRIVER_NAME, dsn)
	if err != nil {
		return nil, fmt.Errorf("open pool: %w", err)
	}

	if err := pool.Ping(); err != nil {
		return nil, fmt.Errorf("ping pool: %w", err)
	}

	return pool, nil
}
