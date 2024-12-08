package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/delivery/http/middleware"
	"github.com/itsLeonB/time-tracker/internal/delivery/http/route"
	"github.com/itsLeonB/time-tracker/internal/provider"
)

type App struct {
	Router *gin.Engine
	config *config
}

type config struct {
	env  string
	port string
}

func (a *App) loadConfig() {
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

	a.config = &config{env, port}
}

func SetupApp() *App {
	db := NewGormDB()
	repositories := provider.ProvideRepositories(db)
	services := provider.ProvideServices(repositories)
	handlers := provider.ProvideHandlers(services)

	a := new(App)
	a.loadConfig()

	gin.SetMode(a.config.env)
	r := gin.Default()
	r.Use(middleware.HandleError())
	r = route.SetupRoutes(r, handlers)

	a.Router = r

	return a
}

func (a *App) Serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.port),
		Handler: a.Router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("error server listen and serve: %s", err.Error())
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("error shutting down: %e", err)
	}

	log.Println("server successfully shutdown")
}
