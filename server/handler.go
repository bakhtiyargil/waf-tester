package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"waf-tester/client"
	"waf-tester/logger"
	"waf-tester/model"
	"waf-tester/service"
)

type Handler struct {
	service *service.TesterService
	logger  *logger.AppLogger
}

func NewHandler() *Handler {
	return &Handler{service: service.NewTesterService(client.NewClient())}
func NewHandler(wp utility.Worker, logger *logger.AppLogger) *Handler {
	return &Handler{
		service: service.NewTesterService(client.NewClient(), wp, logger),
		logger:  logger,
	}
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
			h.logger.Error(err)
		}
		err := h.service.StartInjectionTest(requestBody)
		if err != nil {
			h.logger.Error(err)
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse())

		}
		return c.JSON(http.StatusOK, model.SuccessResponse())
	})
}
