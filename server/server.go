package server

import (
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
	echo *echo.Echo
	cfg  *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{echo: echo.New(), cfg: cfg}
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
	s.AppendMiddlewareHandlers(s.echo)
	s.AppendRouteHandlers(s.echo)

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
