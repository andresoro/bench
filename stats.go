package main

import "time"

// Stats holds individual statistics from requests
type Stats struct {
	ResponseSize  int
	ResponseDur   time.Duration
	TotalRequests uint32
}
