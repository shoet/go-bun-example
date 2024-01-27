package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("failed to load .env: %v\n", err)
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
