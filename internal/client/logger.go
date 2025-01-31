package client

import (
	"log"
	"net/http"
)

// LoggingRoundTripper custom round tripper for logging requests and responses
type LoggingRoundTripper struct{}

func (l LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("[%p] %s %s", req, req.Method, req.URL)

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Println(err)
		return resp, err
	}

	log.Printf("[%p] %d %s", resp.Request, resp.StatusCode, resp.Request.URL)

	return resp, err
}
