package config

import (
	"fmt"
	"os"
	"time"
)

type Auth struct {
	SecretKey      string
	TokenDuration  time.Duration
	CookieDuration time.Duration
	Issuer         string
	URL            string
}

func loadAuthConfig() *Auth {
	return &Auth{
		SecretKey:      os.Getenv("SECRET_KEY"),
		TokenDuration:  5 * time.Hour,
		CookieDuration: 24 * time.Hour,
		Issuer:         os.Getenv("APP_NAME"),
		URL:            fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")),
	}
}
