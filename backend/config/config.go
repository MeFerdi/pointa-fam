package config

import (
	"os"
)

type Config struct {
	DBPath string
}

func LoadConfig() Config {
	return Config{
		DBPath: getEnv("DB_PATH", "./data/pointafam.db"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
