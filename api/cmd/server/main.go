package main

import (
	"caloriecounter/api/internal/ai"
	"caloriecounter/api/internal/auth"
	"caloriecounter/api/internal/config"
	"caloriecounter/api/internal/database"
	"caloriecounter/api/internal/food"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		return
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		slog.Error("database ping failed", "error", err)
		return
	}
	if err = database.RunMigrations(ctx, pool, "db/migrations"); err != nil {
		slog.Error("database migrations failed", "error", err)
		return
	}
	slog.Info("database migrations completed")
	tokens := auth.TokenService{Secret: []byte(cfg.JWTSecret), Expiration: cfg.JWTExpiration}
	authHandler := auth.NewHandler(auth.NewService(pool, tokens))
	foodHandler := food.NewHandler(food.NewService(food.NewRepository(pool), ai.Provider{}, 0.65))
	e := echo.New()
	e.Use(requestLoggingMiddleware)
	e.GET("/health", func(c *echo.Context) error {
		if err := pool.Ping(c.Request().Context()); err != nil {
			slog.Error("health check failed", "error", err)
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
		}
		slog.Debug("health check passed")
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})
	api := e.Group("/api/v1")
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	protected := api.Group("", tokens.Middleware)
	protected.GET("/me", authHandler.Me)
	protected.GET("/me/recommendations", authHandler.Recommendations)
	protected.PUT("/me/profile", authHandler.UpdateProfile)
	protected.POST("/foods/search", foodHandler.Search)
	go func() {
		if serveErr := (echo.StartConfig{Address: ":" + cfg.ServerPort, HideBanner: true}).Start(ctx, e); serveErr != nil && !errors.Is(serveErr, context.Canceled) {
			slog.Error("server stopped", "error", serveErr)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	<-shutdownCtx.Done()
}

func requestLoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		started := time.Now()
		err := next(c)
		status := http.StatusOK
		if httpError, ok := err.(*echo.HTTPError); ok {
			status = httpError.Code
		}
		attrs := []any{"method", c.Request().Method, "path", c.Request().URL.Path, "status", status, "duration_ms", time.Since(started).Milliseconds()}
		if err != nil {
			slog.Error("http request completed with error", append(attrs, "error", err)...)
		} else {
			slog.Info("http request completed", attrs...)
		}
		return err
	}
}
