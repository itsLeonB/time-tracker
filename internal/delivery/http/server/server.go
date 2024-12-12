package server

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
	"github.com/itsLeonB/time-tracker/internal/config"
	"github.com/itsLeonB/time-tracker/internal/delivery/http/middleware"
	"github.com/itsLeonB/time-tracker/internal/delivery/http/route"
	"github.com/itsLeonB/time-tracker/internal/provider"
)

type Server struct {
	Router *gin.Engine
	Config *config.Config
}

func SetupServer() *Server {
	configs := config.LoadConfig()

	db := config.NewGormDB(configs.DB)
	repositories := provider.ProvideRepositories(db)
	services := provider.ProvideServices(configs, repositories)
	handlers := provider.ProvideHandlers(services)

	gin.SetMode(configs.App.Env)
	r := gin.Default()
	r.Use(middleware.HandleError())
	route.SetupRoutes(r, handlers)

	return &Server{
		Router: r,
		Config: configs,
	}
}

func (a *Server) Serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", a.Config.App.Port),
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
