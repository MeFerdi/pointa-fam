package config

import (
	"os"
)

type Config struct {
	DBPath string // Path to the SQLite database file
}

func LoadConfig() *Config {
	return &Config{
		DBPath: os.Getenv("DB_PATH"), // Set this in your .env file
	}
}
