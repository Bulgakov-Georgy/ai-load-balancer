package main

import (
	"github.com/labstack/echo/v4/middleware"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// SimulateProcessingTime adds random response latency
func SimulateProcessingTime(min, max int) {
	delay := time.Duration(rand.Intn(max-min)+min) * time.Millisecond
	time.Sleep(delay)
}

// MockOpenAI simulates OpenAI API response
func MockOpenAI(c echo.Context) error {
	SimulateProcessingTime(1000, 2000)
	return c.String(http.StatusOK, "{\"choices\":[{\"message\":{\"content\":\"OpenRouterResponse\"}}]}")
}

// MockOpenRouter simulates OpenRouter API response
func MockOpenRouter(c echo.Context) error {
	SimulateProcessingTime(1500, 2500)
	return c.String(http.StatusOK, "{\"choices\":[{\"text\":\"OpenAIResponse\"}]}")
}

// MockHuggingFace simulates HuggingFace API response
func MockHuggingFace(c echo.Context) error {
	SimulateProcessingTime(2000, 3000)
	return c.String(http.StatusOK, "[{\"generated_text\":\"HuggingFaceResponse\"}]")
}

func main() {
	e := echo.New()

	// Mock endpoints
	e.Use(middleware.Logger())
	e.POST("/openai", MockOpenAI)
	e.POST("/openrouter", MockOpenRouter)
	e.POST("/huggingface/*", MockHuggingFace)

	log.Println("Starting mock server on :8081")
	e.Logger.Fatal(e.Start(":8081"))
}
