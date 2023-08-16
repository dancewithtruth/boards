package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

const (
	keyAmqpHost     = "AMQP_HOST"
	keyAmqpPort     = "AMQP_PORT"
	keyAmqpUser     = "AMQP_USER"
	keyAmqpPassword = "AMQP_PASSWORD"

	keyEnv    = "ENV"
	valEnvDev = "DEVELOPMENT"
)

// Config encapsulates all the server configuration values.
type Config struct {
	Amqp AmqpConfig
}

// Load looks for config values in environment table and .env files (development), and sets them
// into the Config struct.
func Load(file string) (*Config, error) {
	env := os.Getenv(keyEnv)
	if env == valEnvDev {
		// Load .env file if in development
		err := godotenv.Load(file)
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	amqpConfig, err := getAmqpConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Amqp: amqpConfig,
	}, nil
}

// AmqpConfig represents the config for connecting to a message broker
type AmqpConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
}

// Validate checks that all values are properly loaded into the redis config.
func (config *AmqpConfig) Validate() error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("missing amqp env var: %v", err)
	}
	return nil
}

// getAmqpConfig looks for amqp env vars and creates a config.
func getAmqpConfig() (AmqpConfig, error) {
	cfg := AmqpConfig{
		Host:     os.Getenv(keyAmqpHost),
		Port:     os.Getenv(keyAmqpPort),
		User:     os.Getenv(keyAmqpUser),
		Password: os.Getenv(keyAmqpPassword),
	}

	// validate all redis params are available
	if err := cfg.Validate(); err != nil {
		return AmqpConfig{}, err
	}

	return cfg, nil
}
