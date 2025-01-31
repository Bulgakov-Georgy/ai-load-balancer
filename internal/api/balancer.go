package api

import (
	"errors"
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"ai_load_balancer/internal/cache"
	"ai_load_balancer/internal/client"
	"ai_load_balancer/internal/configuration"
	"ai_load_balancer/internal/metrics"
)

// FetchResponse routes the request to the best available API
func FetchResponse(model string, prompt string) (string, error) {
	responseCacheKey := model + prompt
	if cachedResponse, err := cache.GetString(responseCacheKey); err == nil && cachedResponse != "" {
		return cachedResponse, nil
	}
	apis, ok := client.ModelToApis[model]
	if !ok {
		return "", errors.New("model not found")
	}

	// If we want to implement smart request routing based on the current RPM and API response times we can do something like this
	/*
		sort.Slice(apis, func(i, j int) bool {
			firstRpm := cache.GetRPM(apis[i].Name)
			secondRpm := cache.GetRPM(apis[j].Name)
			if firstRpm == secondRpm {
				return metrics.GetAverageResponseTime(apis[i].Name) < metrics.GetAverageResponseTime(apis[j].Name)
			}
			return firstRpm < secondRpm
		})
	*/

	// This implementation is a simple round-robin which satisfies fallback mechanism described in requirements
	// Basically we'll just prioritize sending requests to the first API in the array that didn't reach limit yet
	for _, api := range apis {
		if rpm, err := cache.GetInt("rpm-" + api.Name); err == nil && rpm < api.MaxRpm {
			metrics.RequestCounter.WithLabelValues(api.Name).Inc()
			cache.Increment("rpm-"+api.Name, time.Minute)
			if response, err := api.Execute(api.Model, prompt); err == nil {
				cache.Set(responseCacheKey, response, configuration.Get().CacheResponseExpirationTime)
				return response, nil
			} else {
				log.Println(err)
			}
		}
	}

	return "", echo.ErrTooManyRequests
}
