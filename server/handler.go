package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"waf-tester/domain/model"
	"waf-tester/logger"
	"waf-tester/service"
)

type Handler interface {
	mapHealthRouteHandlers(health *echo.Group)
	mapBaseRouteHandlers(base *echo.Group)
}

type InjectionTestHandler struct {
	tester service.Tester
	logger logger.Logger
}

func NewInjectionTestHandler(tester service.Tester, logger logger.Logger) Handler {
	return &InjectionTestHandler{
		tester: tester,
		logger: logger,
	}
}

func (h *InjectionTestHandler) mapHealthRouteHandlers(health *echo.Group) {
	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, model.SuccessResponse())
	})
}

func (h *InjectionTestHandler) mapBaseRouteHandlers(base *echo.Group) {
	base.DELETE("/:id/terminate", func(c echo.Context) error {
		testId := c.Param("id")
		err := h.tester.Terminate(testId)
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
		id, err := h.tester.Start(requestBody)
		if err != nil {
			h.logger.Error(err)
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse())

		}
		return c.JSON(http.StatusOK, model.SuccessResponseWithId(id))
	})
}
