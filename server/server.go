package server

import (
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"waf-tester/config"
	"waf-tester/logger"
)
import "github.com/labstack/echo/v4"

type Server interface {
	Start()
}

type TesterServer struct {
	echo    *echo.Echo
	handler Handler
	cfg     *config.Config
	logger  logger.Logger
}

func NewServer(handler Handler, cfg *config.Config, logger logger.Logger) Server {
	return &TesterServer{
		echo:    echo.New(),
		handler: handler,
		cfg:     cfg,
		logger:  logger,
	}
}

func (s *TesterServer) Start() {
	server := &http.Server{
		Addr:           ":" + s.cfg.Server.Default.Port,
		ReadTimeout:    s.cfg.Server.Default.ReadTimeout * time.Second,
		WriteTimeout:   s.cfg.Server.Default.WriteTimeout * time.Second,
		IdleTimeout:    s.cfg.Server.Default.IdleTimeout * time.Second,
		MaxHeaderBytes: s.cfg.Server.Default.MaxHeaderBytes,
	}

	go func() {
		s.logger.Infof("starting server on port %s", s.cfg.Server.Default.Port)
		if err := s.echo.StartServer(server); err != nil {
			log.Fatalf("error starting server: %v", err)
		}
	}()
	s.appendMiddlewares(s.echo)
	s.appendRoutes(s.echo)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	s.logger.Info("shutting down server")
	err := s.echo.Server.Shutdown(ctx)
	if err != nil {
		s.logger.Error(err)
		return
	}
}

func (s *TesterServer) appendMiddlewares(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderXRequestID},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
}

func (s *TesterServer) appendRoutes(e *echo.Echo) {
	base := e.Group("/tests")
	s.handler.mapBaseRouteHandlers(base)

	health := base.Group("/health")
	s.handler.mapHealthRouteHandlers(health)
}
