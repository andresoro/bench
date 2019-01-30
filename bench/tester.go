package bench

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// LoadTester will run repeated requests on an endpoint and aggregate statistics
type LoadTester struct {
	endpoint string
	conns    int
	request  *http.Request
	client   *http.Client
	stats    *Stats
	ch       chan *Stats
	dur      time.Duration
}

// NewTester returns a tester
func NewTester(r *http.Request, conns int, dur time.Duration, end string) *LoadTester {
	return &LoadTester{
		endpoint: end,
		request:  r,
		client:   &http.Client{},
		conns:    conns,
		dur:      dur,
		ch:       make(chan *Stats),
		stats:    &Stats{},
	}
}

// Run initializes the LoadTester with its # of conns for a given duration
// passes the results to the statistics channel
func (l *LoadTester) Run(ch chan *Stats) {
	var wg sync.WaitGroup

	// run a tests for a given duration for all connections
	for i := 0; i < l.conns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for start := time.Now(); time.Since(start) < l.dur; {
				l.test()
			}
		}()
	}
	wg.Wait()

	// average and send the statistics to the upstream channel when finished
	l.stats.avg()
	ch <- l.stats
}

// Make an individual request and update the stats
func (l *LoadTester) test() {
	var body []byte

	start := time.Now()
	resp, err := l.client.Do(l.request)
	if err != nil {
		l.stats.update(0, 0, true)
	}
	rd := time.Since(start)

	if resp != nil {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			l.stats.update(0, 0, true)
		}
	}

	rs := len(body)

	// update stats
	l.stats.update(rs, rd, false)

}
