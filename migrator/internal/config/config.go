package config

type OMSMigratorConfig struct {
	LoggerConfig   loggerConfig `yaml:"logger"`
	PostgresConfig databaseConfig
	MigrationsPath string `env:"MIGRATIONS_PATH"`
}

type loggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type databaseConfig struct {
	Port     string `env:"POSTGRES_PORT"`
	Host     string `env:"POSTGRES_HOST"`
	DBName   string `env:"POSTGRES_DB"`
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	SSLMode  string `env:"POSTGRES_SSL_MODE"`
}
