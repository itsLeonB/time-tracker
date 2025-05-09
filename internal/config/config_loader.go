package config

import (
	"log"
	"os"
)

type Config struct {
	App      *App
	Auth     *Auth
	DB       *DB
	Youtrack *Youtrack
}

func LoadConfig() *Config {
	return &Config{
		App:      loadAppConfig(),
		Auth:     loadAuthConfig(),
		DB:       loadDBConfig(),
		Youtrack: loadYoutrackConfig(),
	}
}

func getRequiredEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s is not set", key)
	}

	return val
}

func getEnvWithDefault(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("%s is not set, using default value: %s", key, defaultVal)
		return defaultVal
	}

	return val
}
