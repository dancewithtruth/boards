package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

const (
	DBHostKey     = "DB_HOST"
	DBPortKey     = "DB_PORT"
	DBNameKey     = "DB_NAME"
	DBUserKey     = "DB_USER"
	DBPasswordKey = "DB_PASSWORD"

	JWTSecretKey     = "JWT_SIGNING_KEY"
	JWTExpirationKey = "JWT_EXPIRATION"

	DockerKey = "DOCKER"
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
	JwtSecret      string
	JwtExpiration  int
}

func Load(file string) (*Config, error) {
	err := godotenv.Load(file)
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	databaseConfig, err := getDatabaseConfig()
	if err != nil {
		return nil, err
	}

	jwtSecret := os.Getenv(JWTSecretKey)
	jwtExpirationStr := os.Getenv(JWTExpirationKey)
	jwtExpiration, err := strconv.Atoi(jwtExpirationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expiration value: %w", err)
	}

	return &Config{
		DatabaseConfig: databaseConfig,
		JwtSecret:      jwtSecret,
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

	// This allows running tests from outside the docker network assuming your local
	// development environment has ports exposed
	if os.Getenv(DockerKey) == "" {
		databaseConfig.Host = "localhost"
	}

	// validate all db params are available
	if err := databaseConfig.Validate(); err != nil {
		return DatabaseConfig{}, err
	}

	return databaseConfig, nil
}
