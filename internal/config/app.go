package config

import (
	"context"
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
}

func SetupApp() *App {
	db := NewGormDB()
	repositories := provider.ProvideRepositories(db)
	services := provider.ProvideServices(repositories)
	handlers := provider.ProvideHandlers(services)

	gin.SetMode(os.Getenv("APP_ENV"))
	r := gin.Default()
	r.Use(middleware.HandleError())
	r = route.SetupRoutes(r, handlers)

	return &App{r}
}

func (a *App) Serve() {
	srv := http.Server{
		Addr:    ":8080", //+ os.Getenv("APP_PORT"),
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
