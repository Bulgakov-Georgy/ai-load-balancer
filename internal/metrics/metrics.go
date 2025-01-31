package metrics

import "github.com/prometheus/client_golang/prometheus"

// RequestCounter tracks total API requests
var RequestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of requests",
	},
	[]string{"service"},
)

// APIResponseTime tracks response times
var APIResponseTime = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "api_response_time",
		Help:    "Response times of API calls",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"service"},
)
