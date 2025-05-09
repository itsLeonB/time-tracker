package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/config"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/delivery/http/middleware"
	strategy "github.com/itsLeonB/catfeinated-time-tracker/internal/delivery/http/middleware/strategy/error"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/delivery/http/route"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/provider"
)

type Server struct {
	Router *gin.Engine
	config *config.Config
}

func SetupServer() *Server {
	configs := config.LoadConfig()

	db := config.NewGormDB(configs.DB)
	repositories := provider.ProvideRepositories(db, configs)
	services := provider.ProvideServices(configs, repositories)
	handlers := provider.ProvideHandlers(services)

	errorStrategyMap := strategy.NewErrorStrategyMap()

	gin.SetMode(configs.App.Env)
	r := gin.Default()
	r.Use(middleware.HandleError(errorStrategyMap))
	route.SetupRoutes(r, handlers, services)

	return &Server{
		Router: r,
		config: configs,
	}
}

func (s *Server) Serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", s.config.App.Port),
		Handler: s.Router,
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
		s.config.App.Timeout,
	)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("error shutting down: %s", err.Error())
	}

	log.Println("server successfully shutdown")
}
