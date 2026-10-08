package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct{ service *Service }
type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(service *Service) *Handler        { return &Handler{service: service} }
func (h *Handler) Register(c *echo.Context) error { return h.handle(c, h.service.Register) }
func (h *Handler) Login(c *echo.Context) error    { return h.handle(c, h.service.Login) }
func (h *Handler) handle(c *echo.Context, operation func(context.Context, string, string) (string, error)) error {
	var input credentials
	if err := c.Bind(&input); err != nil {
		slog.Warn("authentication request body could not be parsed", "path", c.Request().URL.Path, "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	slog.Info("authentication request received", "path", c.Request().URL.Path, "email", input.Email)
	token, err := operation(c.Request().Context(), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			slog.Warn("authentication failed", "path", c.Request().URL.Path, "email", input.Email)
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}
		slog.Error("authentication request failed", "path", c.Request().URL.Path, "email", input.Email, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "authentication failed")
	}
	slog.Info("authentication succeeded", "path", c.Request().URL.Path, "email", input.Email)
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
