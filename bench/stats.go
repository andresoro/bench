package bench

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

// update the object with incoming data
func (s *Stats) update(rs int, rd time.Duration, err bool) {
	s.TotalRequests++
	if err {
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
