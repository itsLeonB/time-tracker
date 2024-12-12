package config

import (
	"log"
	"os"
)

type Config struct {
	App  *App
	Auth *Auth
	DB   *DB
}

func LoadConfig() *Config {
	return &Config{
		App:  loadAppConfig(),
		Auth: loadAuthConfig(),
		DB:   loadDBConfig(),
	}
}

func getRequiredEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s is not set", key)
	}

	return val
}
