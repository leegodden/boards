package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

const (
	DBHostKey     = "DB_HOST"
	DBPortKey     = "DB_PORT"
	DBNameKey     = "DB_NAME"
	DBUserKey     = "DB_USER"
	DBPasswordKey = "DB_PASSWORD"
)

type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Name     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
}

// Ensure all required fields in the DatabaseConfig are properly set
func (dbConfig *DatabaseConfig) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dbConfig); err != nil {
		return fmt.Errorf("missing database env var: %v", err)
	}
	return nil
}

// Extensible: can easily be expanded to include more configurations
type Config struct {
	DatabaseConfig DatabaseConfig
}

// Manages the fetching and validation of configuration data
func Load() (*Config, error) {
	databaseConfig, err := getDatabaseConfig()
	if err != nil {
		return nil, err
	}
	return &Config{DatabaseConfig: databaseConfig}, nil
}

// Get the env variables for the database connection
func getDatabaseConfig() (DatabaseConfig, error) {
	databaseConfig := DatabaseConfig {
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