package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

func (s TokenService) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		parts := strings.SplitN(c.Request().Header.Get("Authorization"), " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			slog.Warn("request rejected: missing bearer token", "method", c.Request().Method, "path", c.Request().URL.Path)
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		userID, err := s.Parse(parts[1])
		if err != nil {
			slog.Warn("request rejected: invalid bearer token", "method", c.Request().Method, "path", c.Request().URL.Path, "error", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		slog.Debug("authenticated request", "method", c.Request().Method, "path", c.Request().URL.Path, "user_id", userID)
		c.Set("user_id", userID)
		return next(c)
	}
}
