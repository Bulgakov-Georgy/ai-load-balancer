package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"ai_load_balancer/internal/configuration"
	"ai_load_balancer/internal/metrics"
)

type openRouterRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type openRouterResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type openRouterError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

var openRouterClient = &http.Client{Transport: LoggingRoundTripper{}, Timeout: configuration.Get().OpenRouterTimeout}

func fetchOpenRouterResponse(model string, prompt string) (string, error) {
	start := time.Now()
	requestBody, err := json.Marshal(openRouterRequest{model, prompt})
	if err != nil {
		return "", err
	}

	config := configuration.Get()
	req, err := http.NewRequest("POST", config.OpenRouterUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+config.OpenRouterKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := openRouterClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	metrics.APIResponseTime.WithLabelValues("openrouter").Observe(time.Since(start).Seconds())

	var result openRouterResponse
	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {
		err = json.Unmarshal(byteBody, &result)
	} else {
		var responseError openRouterError
		err = json.Unmarshal(byteBody, &responseError)
		if err == nil {
			return "", errors.New(responseError.Error.Message)
		}
		return "", err
	}

	text := result.Choices[0].Text
	if text == "" {
		return "", errors.New("couldn't find OpenRouter response")
	}
	return text, nil
}
