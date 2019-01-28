package main

import (
	"time"
)

// Stats holds individual statistics from requests
type Stats struct {
	Endpoint      string
	ResponseSize  float64
	ResponseDur   time.Duration
	TotalRequests int64
	err           int64
}
