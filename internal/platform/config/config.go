package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	// db
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
		//
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),
	}

	if config.JWTSecret == "" {
		return nil, errors.New("missing JWT_SECRET from .env")
	}

	if config.DBHost == "" {
		config.DBHost = "localhost"
	}
	if config.DBPort == "" {
		config.DBPort = "5432"
	}
	if config.DBUser == "" {
		return nil, errors.New("missing DB_USER from .env")
	}
	if config.DBName == "" {
		return nil, errors.New("missing DB_NAME from .env")
	}
	if config.DBPass == "" {
		return nil, errors.New("missing DB_PASS from .env")
	}
	return &config, nil
}
