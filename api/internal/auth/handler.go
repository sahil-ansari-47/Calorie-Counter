package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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
func (h *Handler) Me(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	profile, err := h.service.GetProfile(c.Request().Context(), userID)
	if err != nil {
		slog.Error("profile lookup failed", "user_id", userID, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "profile lookup failed")
	}
	return c.JSON(http.StatusOK, map[string]any{"profile": profile})
}
func (h *Handler) Recommendations(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	profile, err := h.service.GetRecommendations(c.Request().Context(), userID)
	if err != nil {
		slog.Error("recommendations lookup failed", "user_id", userID, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "recommendations lookup failed")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"maintenance_calories":            profile.MaintenanceCalories,
		"recommended_calories":            profile.RecommendedCalories,
		"recommended_protein_g_day":       profile.RecommendedProteinGDay,
		"recommended_carbohydrates_g_day": profile.RecommendedCarbohydratesGDay,
		"recommended_fat_g_day":           profile.RecommendedFatGDay,
		"recommended_fiber_g_day":         profile.RecommendedFiberGDay,
		"recommended_sugars_g_day":        profile.RecommendedSugarsGDay,
		"recommended_saturated_fat_g_day": profile.RecommendedSaturatedFatGDay,
	})
}
func (h *Handler) UpdateProfile(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	var input ProfileUpdate
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	profile, err := h.service.UpdateProfile(c.Request().Context(), userID, input)
	if err != nil {
		slog.Error("profile update failed", "user_id", userID, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "profile update failed")
	}
	return c.JSON(http.StatusOK, map[string]any{"profile": profile})
}
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
