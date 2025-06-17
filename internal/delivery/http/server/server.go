package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/delivery/http/route"
	"github.com/itsLeonB/time-tracker/internal/provider"
)

func SetupHTTPServer(configs *ezutil.Config) *http.Server {
	repositories := provider.ProvideRepositories(configs)
	services := provider.ProvideServices(configs, repositories)
	handlers := provider.ProvideHandlers(services)

	gin.SetMode(configs.App.Env)
	r := gin.Default()
	route.SetupRoutes(r, handlers, services)

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", configs.App.Port),
		Handler: r,
	}
}

func DefaultConfigs() ezutil.Config {
	timeout, _ := time.ParseDuration("10s")
	tokenDuration, _ := time.ParseDuration("24h")
	cookieDuration, _ := time.ParseDuration("24h")
	secretKey, err := ezutil.GenerateRandomString(32)
	if err != nil {
		log.Fatal("error generating secret key: %w", err)
	}

	appConfig := ezutil.App{
		Env:        "debug",
		Port:       "8080",
		Timeout:    timeout,
		ClientUrls: []string{"http://localhost:3000"},
		Timezone:   "Asia/Jakarta",
	}

	authConfig := ezutil.Auth{
		SecretKey:      secretKey,
		TokenDuration:  tokenDuration,
		CookieDuration: cookieDuration,
		Issuer:         "time-tracker",
		URL:            "http://localhost:8000",
	}

	return ezutil.Config{
		App:  &appConfig,
		Auth: &authConfig,
	}
}
