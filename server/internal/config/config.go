package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
)

const (
	DBHostKey     = "DB_HOST"
	DBPortKey     = "DB_PORT"
	DBNameKey     = "DB_NAME"
	DBUserKey     = "DB_USER"
	DBPasswordKey = "DB_PASSWORD"

	JWTSigningKey = "JWT_SIGNING_KEY"
	JWTExpiration = "JWT_EXPIRATION"
)

type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Name     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
}

func (dbConfig *DatabaseConfig) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dbConfig); err != nil {
		return fmt.Errorf("missing database env var: %v", err)
	}
	return nil
}

type Config struct {
	DatabaseConfig DatabaseConfig
	JwtSigningKey  string
	JwtExpiration  int
}

func Load() (*Config, error) {
	databaseConfig, err := getDatabaseConfig()
	if err != nil {
		return nil, err
	}

	jwtSigningKey := os.Getenv(JWTSigningKey)
	jwtExpirationStr := os.Getenv(JWTExpiration)
	jwtExpiration, err := strconv.Atoi(jwtExpirationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expiration value: %w", err)
	}

	return &Config{
		DatabaseConfig: databaseConfig,
		JwtSigningKey:  jwtSigningKey,
		JwtExpiration:  jwtExpiration,
	}, nil
}

func getDatabaseConfig() (DatabaseConfig, error) {
	databaseConfig := DatabaseConfig{
		Host:     os.Getenv(DBHostKey),
		Port:     os.Getenv(DBPortKey),
		Name:     os.Getenv(DBNameKey),
		User:     os.Getenv(DBUserKey),
		Password: os.Getenv(DBPasswordKey),
	}

	// validate all db params are available
	if err := databaseConfig.Validate(); err != nil {
		return DatabaseConfig{}, err
	}

	return databaseConfig, nil
}
