package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

const (
	keyDBHost     = "DB_HOST"
	keyDBPort     = "DB_PORT"
	keyDBName     = "DB_NAME"
	keyDBUser     = "DB_USER"
	keyDBPassword = "DB_PASSWORD"

	keyEnv             = "ENV"
	keyServerPort      = "SERVER_PORT"
	keyJWTSecret       = "JWT_SIGNING_KEY"
	keyJWTExpiration   = "JWT_EXPIRATION"
	keyInternalNetwork = "INTERNAL_NETWORK"

	valEnvDev = "DEVELOPMENT"
)

// Config encapsulates all the server configuration values.
type Config struct {
	ServerPort    string
	JwtSecret     string
	JwtExpiration int
	DB            DatabaseConfig
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

	databaseConfig, err := getDatabaseConfig()
	if err != nil {
		return nil, err
	}

	serverPort := os.Getenv(keyServerPort)
	jwtSecret := os.Getenv(keyJWTSecret)
	jwtExpirationStr := os.Getenv(keyJWTExpiration)

	jwtExpiration, err := strconv.Atoi(jwtExpirationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expiration value: %w", err)
	}

	return &Config{
		DB:            databaseConfig,
		ServerPort:    serverPort,
		JwtSecret:     jwtSecret,
		JwtExpiration: jwtExpiration,
	}, nil
}

// DatabaseConfig encapsulates all the config values for the database.
type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Name     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
}

// Validate checks that all values are properly loaded into the database config.
func (dbConfig *DatabaseConfig) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dbConfig); err != nil {
		return fmt.Errorf("missing database env var: %v", err)
	}
	return nil
}

func getDatabaseConfig() (DatabaseConfig, error) {
	databaseConfig := DatabaseConfig{
		Host:     os.Getenv(keyDBHost),
		Port:     os.Getenv(keyDBPort),
		Name:     os.Getenv(keyDBName),
		User:     os.Getenv(keyDBUser),
		Password: os.Getenv(keyDBPassword),
	}

	// This allows running tests from outside the docker network assuming your local
	// development environment has ports exposed
	if os.Getenv(keyInternalNetwork) == "false" {
		databaseConfig.Host = "localhost"
	}

	// validate all db params are available
	if err := databaseConfig.Validate(); err != nil {
		return DatabaseConfig{}, err
	}

	return databaseConfig, nil
}
