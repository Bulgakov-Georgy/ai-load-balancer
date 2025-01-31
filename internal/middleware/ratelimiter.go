package middleware

import (
	"time"

	"github.com/labstack/echo/v4"

	"ai_load_balancer/internal/cache"
	"ai_load_balancer/internal/configuration"
)

// RateLimiterMiddleware middleware for limiting amount of requests per user
type RateLimiterMiddleware struct{}

func (RateLimiterMiddleware) Allow(identifier string) (bool, error) {
	key := "user-rpm-" + identifier

	count, err := cache.GetInt(key)
	if err != nil {
		return false, err
	}

	if count >= configuration.Get().UserRpm {
		return false, echo.ErrTooManyRequests
	}

	cache.Increment(key, time.Minute)
	return true, nil
}
