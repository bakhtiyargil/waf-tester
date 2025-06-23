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

func (s *Server) Start() error {
	server := &http.Server{
		Addr:           s.cfg.Server.Default.Port,
		ReadTimeout:    time.Second * s.cfg.Server.Default.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.Default.WriteTimeout,
		MaxHeaderBytes: s.cfg.Server.Default.MaxHeaderBytes,
	}

	go func() {
		log.Println("Server is starting at http://localhost" + s.cfg.Server.Default.Port)
		if err := s.echo.StartServer(server); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	//add handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	log.Println("Server exited properly")
	return s.echo.Server.Shutdown(ctx)
}
