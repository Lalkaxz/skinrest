package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Host    string `envconfig:"SERVER_HOST" default:"localhost"`
	Port    int    `envconfig:"SERVER_PORT" default:"8081"`
	ApiEnv  string `envconfig:"API_ENV" default:"local"`
	GinMode string `envconfig:"GIN_MODE" default:"debug"`
}

type DatabaseConfig struct {
	Driver   string `envconfig:"DATABASE_DRIVER" default:"postgres"`
	Host     string `envconfig:"DATABASE_HOST" default:"localhost"`
	Port     int    `envconfig:"DATABASE_PORT" default:"5432"`
	User     string `envconfig:"DATABASE_USER" required:"true"`
	Password string `envconfig:"DATABASE_PASSWORD" required:"true"`
	Name     string `envconfig:"DATABASE_NAME" required:"true"`
	SSL      string `envconfig:"DATABASE_SSL" default:"disable"`
}

type AuthConfig struct {
	JwtSecret string `envconfig:"AUTH_JWT_SECRET" required:"true"`
}

func GetConfig() *Config {
	var config Config

	// load values from environment variables
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Failed to load config: %v ", err)
	}

	return &config
}
