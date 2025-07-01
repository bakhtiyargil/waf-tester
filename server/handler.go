package server

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"waf-tester/client"
	"waf-tester/model"
	"waf-tester/service"
)

type Handler struct {
	service *service.TesterService
}

func NewHandler() *Handler {
	return &Handler{service: service.NewTesterService(client.NewClient())}
}

func (h *Handler) mapHealthRouteHandlers(health *echo.Group) {
	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, model.SuccessResponse())
	})
}

func (h *Handler) mapBaseRouteHandlers(base *echo.Group) {
	base.POST("", func(c echo.Context) error {
		requestBody := new(model.TestRequest)
		if err := c.Bind(requestBody); err != nil {
			log.Panicf("error binding request body %v", err)
		}
		err := h.service.StartInjectionTest(requestBody)
		if err != nil {
			log.Printf("unexpected internal error: %v", err)
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse())

		}
		return c.JSON(http.StatusOK, model.SuccessResponse())
	})
}
