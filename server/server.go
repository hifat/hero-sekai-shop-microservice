package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middleware"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middleware/middlewareHandler"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	app        *echo.Echo
	db         *mongo.Client
	cfg        *config.Config
	middleware middlewareHandler.Handler
}

func (s *server) gracefulShutdown(pctx context.Context, quit <-chan os.Signal) {
	logger.Info(fmt.Sprintf("start service: %s", s.cfg.App.Name))
	<-quit
	logger.Info(fmt.Sprintf("shuting down service: %s", s.cfg.App.Name))

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	if err := s.app.Shutdown(ctx); err != nil {
		logger.Fatal(err.Error())
	}
}

func (s *server) httpListening() {
	if err := s.app.Start(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err.Error())
	}
}

func Start(pctx context.Context, cfg *config.Config, db *mongo.Client) {
	s := &server{
		app:        echo.New(),
		db:         db,
		cfg:        cfg,
		middleware: middleware.InitMiddleware(cfg, db),
	}

	s.app.Use(echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Skipper:      echoMiddleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      30 * time.Second,
	}))

	s.app.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		Skipper:      echoMiddleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PATCH, echo.DELETE},
	}))

	s.app.Use(echoMiddleware.BodyLimit("10M"))

	switch s.cfg.App.Name {
	case "auth":
		s.authService()
	case "player":
		s.playerService()
	case "item":
		s.itemService()
	case "inventory":
		s.inventoryService()
	case "payment":
		s.paymentService()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s.app.Use(echoMiddleware.Logger())

	go s.gracefulShutdown(pctx, quit)

	s.httpListening()
}
