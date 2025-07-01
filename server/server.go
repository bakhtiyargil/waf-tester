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
)
import "github.com/labstack/echo/v4"

type Server struct {
	echo    *echo.Echo
	cfg     *config.Config
	handler *Handler
}

func NewServer(cfg *config.Config, handler *Handler) *Server {
	return &Server{echo: echo.New(), cfg: cfg, handler: handler}
}

func (s *Server) Start() {
	server := &http.Server{
		Addr:           ":" + s.cfg.Server.Default.Port,
		ReadTimeout:    s.cfg.Server.Default.ReadTimeout * time.Second,
		WriteTimeout:   s.cfg.Server.Default.WriteTimeout * time.Second,
		IdleTimeout:    s.cfg.Server.Default.IdleTimeout * time.Second,
		MaxHeaderBytes: s.cfg.Server.Default.MaxHeaderBytes,
	}

	go func() {
		log.Println("server is starting at http://localhost:" + s.cfg.Server.Default.Port)
		if err := s.echo.StartServer(server); err != nil {
			log.Fatalf("error starting server: %v", err)
		}
	}()
	s.AppendMiddlewares(s.echo)
	s.AppendRoutes(s.echo)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	log.Println("server exited properly")
	err := s.echo.Server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("error shutting down server: %v", err)
	}
}

func (s *Server) AppendMiddlewares(e *echo.Echo) {
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

func (s *Server) AppendRoutes(e *echo.Echo) {
	base := e.Group("/test")
	s.handler.mapBaseRouteHandlers(base)

	health := base.Group("/health")
	s.handler.mapHealthRouteHandlers(health)
}
