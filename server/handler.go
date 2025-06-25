package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"waf-tester/client"
	"waf-tester/model"
	"waf-tester/service"
)

func (s *Server) AppendMiddlewareHandlers(e *echo.Echo) {
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

func (s *Server) AppendRouteHandlers(e *echo.Echo) {
	base := e.Group("/test")
	mapBaseRouteHandlers(base)

	health := base.Group("/health")
	mapHealthRouteHandlers(health)
}

func mapHealthRouteHandlers(health *echo.Group) {
	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, model.SuccessResponse())
	})
}

func mapBaseRouteHandlers(base *echo.Group) {
	base.POST("", func(c echo.Context) error {
		requestBody := new(model.TestRequest)
		if err := c.Bind(requestBody); err != nil {
			log.Panicf("error binding request body %v", err)
		}
		svc := service.NewTesterService(&client.Client{})
		result, err := svc.StartInjectionTest(requestBody)
		if err != nil {
			log.Panicf("error while starting test %v", err)
		}
		return c.JSON(http.StatusOK, result)
	})
}
