package config

type OMSNotifyServiceConfig struct {
	LoggerConfig   loggerConfig `yaml:"logger"`
	PostgresConfig databaseConfig
	KafkaConfig    brokerConfig `yaml:"kafka"`
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

type brokerConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	GroupID string   `yaml:"group_id"`
}
