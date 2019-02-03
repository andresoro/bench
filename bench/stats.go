package bench

import (
	"fmt"
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

// update the object with incoming data
func (s *Stats) update(rs int, rd time.Duration, err error) {
	s.TotalRequests++
	if err != nil {
		s.err++
		return
	}
	s.ResponseSize += float64(rs)
	s.ResponseDur += rd
	return
}

// change responseSize, responseDur to be averages based on total requests
// should be called after the tester is done
// does not count errors towards time or size averages
func (s *Stats) avg() {
	// average response duration
	s.ResponseDur = s.ResponseDur / time.Duration(s.TotalRequests-s.err)

	//average response size
	s.ResponseSize = s.ResponseSize / float64(s.TotalRequests-s.err)
}

func (s *Stats) print() {
	fmt.Printf("Test completed for endpoint: %s \n", s.Endpoint)
	fmt.Printf("	Total requests completed: %d \n", s.TotalRequests)
	fmt.Printf("	Total errors: %d \n", s.err)
	fmt.Printf("	Average response size: %f bytes\n", s.ResponseSize)
	fmt.Printf("	Average response time: %s \n", s.ResponseDur.String())
}
