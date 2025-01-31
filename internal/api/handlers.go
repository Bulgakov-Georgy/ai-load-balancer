package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"ai_load_balancer/internal/metrics"
)

// GenerateRequest represents the incoming request payload
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// GenerateResponse represents the API response
type GenerateResponse struct {
	Response string `json:"response"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// HandleGenerate processes the generate request
// @Summary Generate AI Response
// @Description Generates text using one of the available AI APIs
// @Accept json
// @Produce json
// @Param request body GenerateRequest true "prompt request"
// @Success 200 {object} GenerateResponse
// @Failure 403 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /generate [post]
func HandleGenerate(c echo.Context) error {
	metrics.RequestCounter.WithLabelValues("ai_load_balancer").Inc()

	var req GenerateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	response, err := FetchResponse(req.Model, req.Prompt)
	switch {
	case errors.Is(err, echo.ErrTooManyRequests):
		return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "Too many requests"})
	case err == nil:
		return c.JSON(http.StatusOK, GenerateResponse{Response: response})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get response"})
	}
}
