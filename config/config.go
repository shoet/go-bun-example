package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func init() {
	if err := godotenv.Load(); err != nil {
		env, ok := os.LookupEnv("ENV")
		if !ok {
			panic("ENV is not set")
		}
		if env == "local" {
			panic(fmt.Sprintf("failed to load .env: %s", err))
		}
	}
}

type DBConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

func NewDBConfig() (*DBConfig, error) {
	var dbConfig DBConfig
	if err := envconfig.Process("DB", &dbConfig); err != nil {
		return nil, fmt.Errorf("failed to process env var: %w", err)
	}
	return &dbConfig, nil
}
