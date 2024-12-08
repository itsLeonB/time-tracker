package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB() *gorm.DB {
	dsn := getDsn()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting to db: %s", err.Error())
	}

	return db
}

func getDsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		getRequiredEnv("DB_HOST"),
		getRequiredEnv("DB_USER"),
		getRequiredEnv("DB_PASSWORD"),
		getRequiredEnv("DB_NAME"),
		getRequiredEnv("DB_PORT"),
		getRequiredEnv("DB_SSL_MODE"),
		getRequiredEnv("DB_TIME_ZONE"),
	)
}

func getRequiredEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s is not set", key)
	}

	return val
}
