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

type huggingFaceRequest struct {
	Inputs string `json:"inputs"`
}

type huggingFaceResponse []struct {
	GeneratedText string `json:"generated_text"`
}

type huggingFaceError struct {
	Error string `json:"error"`
}

var huggingFaceClient = &http.Client{Transport: LoggingRoundTripper{}, Timeout: configuration.Get().HuggingFaceTimeout}

func fetchHuggingFaceResponse(model string, prompt string) (string, error) {
	start := time.Now()
	requestBody, err := json.Marshal(huggingFaceRequest{Inputs: prompt})
	if err != nil {
		return "", err
	}

	config := configuration.Get()
	req, err := http.NewRequest("POST", config.HuggingFaceUrl+model, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+config.HuggingFaceKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := huggingFaceClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	metrics.APIResponseTime.WithLabelValues("huggingface").Observe(time.Since(start).Seconds())

	var result huggingFaceResponse
	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {
		err = json.Unmarshal(byteBody, &result)
	} else {
		var responseError huggingFaceError
		err = json.Unmarshal(byteBody, &responseError)
		if err == nil {
			return "", errors.New(responseError.Error)
		}
		return "", err
	}

	text := result[0].GeneratedText
	if text == "" {
		return "", errors.New("couldn't find HuggingFace response")
	}
	return text, nil
}
