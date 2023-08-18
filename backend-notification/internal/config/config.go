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

	keyEmailHost     = "EMAIL_HOST"
	keyEmailPort     = "EMAIL_PORT"
	keyEmailFrom     = "EMAIL_FROM"
	keyEmailPassword = "EMAIL_PASSWORD"

	keyBoardsBaseURL = "BOARDS_BASE_URL"

	keyEnv    = "ENV"
	valEnvDev = "DEVELOPMENT"
)

// Config encapsulates all the server configuration values.
type Config struct {
	Amqp          AmqpConfig
	Email         EmailConfig
	BoardsBaseURL string
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

	emailConfig, err := getEmailConfig()
	if err != nil {
		return nil, err
	}

	boardsBaseURL := os.Getenv(keyBoardsBaseURL)

	return &Config{
		Amqp:          amqpConfig,
		Email:         emailConfig,
		BoardsBaseURL: boardsBaseURL,
	}, nil
}

// AmqpConfig represents the config for connecting to a message broker
type AmqpConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
}

// Validate checks that all values are properly loaded into the amqp config.
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

// EmailConfig represents the config for connecting to an smtp server
type EmailConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	From     string `validate:"required"`
	Password string `validate:"required"`
}

// Validate checks that all values are properly loaded into the email config.
func (config *EmailConfig) Validate() error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("missing email env var: %v", err)
	}
	return nil
}

// getEmailConfig looks for email env vars and creates a config.
func getEmailConfig() (EmailConfig, error) {
	cfg := EmailConfig{
		Host:     os.Getenv(keyEmailHost),
		Port:     os.Getenv(keyEmailPort),
		From:     os.Getenv(keyEmailFrom),
		Password: os.Getenv(keyEmailPassword),
	}

	// validate all redis params are available
	if err := cfg.Validate(); err != nil {
		return EmailConfig{}, err
	}

	return cfg, nil
}
