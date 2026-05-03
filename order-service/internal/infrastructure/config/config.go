package config

import "time"

type OMSGOrderServiceConfig struct {
	LoggerConfig     loggerConfig `yaml:"logger"`
	PostgresConfig   databaseConfig
	KafkaConfig      brokerConfig     `yaml:"kafka"`
	DispatcherConfig dispatcherConfig `yaml:"dispatcher"`
	ServerConfig     serverConfig     `yaml:"server"`
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
}

type dispatcherConfig struct {
	ChillDuration time.Duration `yaml:"chillTime"`
	WorkersCount  int           `yaml:"workerCount"`
	BatchSize     int           `yaml:"batchSize"`
	RetryCount    int           `yaml:"retryCount"`
}

type serverConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}
