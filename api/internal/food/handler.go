package food

import (
	"caloriecounter/api/internal/nutrition"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) Search(c *echo.Context) error {
	var request nutrition.Request
	if err := c.Bind(&request); err != nil {
		slog.Warn("food search body could not be parsed", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	slog.Info("food search received", "query", request.Query, "quantity", request.Quantity, "unit", request.Unit)
	result, err := h.service.Search(c.Request().Context(), request)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			slog.Warn("food search validation failed", "query", request.Query, "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
		}
		slog.Error("food search failed", "query", request.Query, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "food search failed")
	}
	slog.Info("food search succeeded", "query", request.Query, "food_id", result.ID)
	return c.JSON(http.StatusOK, map[string]any{"food": result})
}
