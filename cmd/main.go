package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/echo-swagger"

	_ "ai_load_balancer/docs"
	"ai_load_balancer/internal/api"
	"ai_load_balancer/internal/cache"
	"ai_load_balancer/internal/metrics"
	"ai_load_balancer/internal/middleware"
)

// @title AI Load Balancer API
// @version 1.0
// @description This is a load balancer service for AI model APIs.
// @host localhost:8080
// @BasePath /
func main() {
	cache.InitRedis()

	prometheus.MustRegister(metrics.RequestCounter, metrics.APIResponseTime)

	e := echo.New()

	e.Use(echomiddleware.LoggerWithConfig(echomiddleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/metrics") || strings.HasPrefix(c.Path(), "/swagger")
		},
	}))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RateLimiterWithConfig(echomiddleware.RateLimiterConfig{
		Skipper: echomiddleware.DefaultSkipper,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		Store: middleware.RateLimiterMiddleware{},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, api.ErrorResponse{Error: "Couldn't identify user"})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			if errors.Is(err, echo.ErrTooManyRequests) {
				return context.JSON(http.StatusTooManyRequests, api.ErrorResponse{Error: http.StatusText(http.StatusTooManyRequests)})
			}
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, api.ErrorResponse{Error: http.StatusText(http.StatusInternalServerError)})
		},
	}))

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.POST("/generate", api.HandleGenerate)

	log.Println("Starting server on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
