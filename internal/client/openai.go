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

type openAIRequest struct {
	Model    string                 `json:"model"`
	Messages []openAIRequestMessage `json:"messages"`
}

type openAIRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type openAIError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

var openAIClient = &http.Client{Transport: LoggingRoundTripper{}, Timeout: configuration.Get().OpenAITimeout}

func fetchOpenAIResponse(model string, prompt string) (string, error) {
	start := time.Now()
	request := openAIRequest{
		Model:    model,
		Messages: []openAIRequestMessage{{"user", prompt}},
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	config := configuration.Get()
	req, err := http.NewRequest("POST", config.OpenAIUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+config.OpenAIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := openAIClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	metrics.APIResponseTime.WithLabelValues("openai").Observe(time.Since(start).Seconds())

	var result openAIResponse
	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {
		err = json.Unmarshal(byteBody, &result)
	} else {
		var responseError openAIError
		err = json.Unmarshal(byteBody, &responseError)
		if err == nil {
			return "", errors.New(responseError.Error.Message)
		}
		return "", err
	}

	text := result.Choices[0].Message.Content
	if text == "" {
		return "", errors.New("couldn't find OpenAI response")
	}
	return text, nil
}
