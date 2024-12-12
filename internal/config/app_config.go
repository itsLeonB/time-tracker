package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type App struct {
	Env     string
	Port    string
	Timeout time.Duration
}

func loadAppConfig() *App {
	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Println("APP_ENV is not set, using default value: debug")
		env = "debug"
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Println("APP_PORT is not set, using default value: 8000")
		port = "8000"
	}

	timeout := os.Getenv("APP_TIMEOUT_SECONDS")
	if timeout == "" {
		log.Println("APP_TIMEOUT_SECONDS is not set, using default value: 10")
		timeout = "10"
	}

	timeoutInt, err := strconv.Atoi(timeout)
	if err != nil {
		log.Fatalf("error converting APP_TIMEOUT_SECONDS to int: %s", err.Error())
	}

	return &App{
		Env:     env,
		Port:    port,
		Timeout: time.Duration(timeoutInt) * time.Second,
	}
}
