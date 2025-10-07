package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppName       string        `env:"APP_NAME" envDefault:"booking"`
	HTTPAddr      string        `env:"HTTP_ADDR" envDefault:":8080"`
	ShutdownGrace time.Duration `env:"SHUTDOWN_GRACE" envDefault:"10s"`

	//DB
	PostgresDSN string `env:"POSTGRES_DSN"`
	DB          DatabaseConfig

	// Observability
	OTLPEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:"localhost:4317"`
	Env          string `env:"ENV" envDefault:"dev"`
}

type DatabaseConfig struct {
	Host         string        `env:"DB_HOST" envDefault:"localhost"`
	Port         int           `env:"DB_PORT" envDefault:"5432"`
	User         string        `env:"DB_USER"`
	Password     string        `env:"DB_PASSWORD"`
	Name         string        `env:"DB_NAME"`
	SSLMode      string        `env:"DB_SSLMODE" envDefault:"disable"`
	MaxOpenConns int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	MaxLifetime  time.Duration `env:"DB_MAX_LIFETIME" envDefault:"5m"`
}

func (c Config) GetDSN() string {
	if c.PostgresDSN != "" {
		return c.PostgresDSN
	}

	return c.DB.DSN()
}

func (db DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}

func Load() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return cfg, err
}
