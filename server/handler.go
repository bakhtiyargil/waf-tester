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
	logger  logger.Logger
}

func NewHandler(logger logger.Logger) *Handler {
	return &Handler{
		service: service.NewTesterService(client.NewClient(), logger),
		logger:  logger,
	}
}

func (h *Handler) mapHealthRouteHandlers(health *echo.Group) {
	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, model.SuccessResponse())
	})
}

func (h *Handler) mapBaseRouteHandlers(base *echo.Group) {
	base.DELETE("/:id/terminate", func(c echo.Context) error {
		testId := c.Param("id")
		err := h.service.TerminateInjectionTest(testId)
		if err != nil {
			h.logger.Error(err)
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse())

		}
		return c.JSON(http.StatusOK, model.SuccessResponseWithId(testId))
	})

	base.POST("/start", func(c echo.Context) error {
		requestBody := new(model.TestRequest)
		if err := c.Bind(requestBody); err != nil {
			h.logger.Error(err)
		}
		id, err := h.service.StartInjectionTest(requestBody)
		if err != nil {
			h.logger.Error(err)
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse())

		}
		return c.JSON(http.StatusOK, model.SuccessResponseWithId(id))
	})
}
