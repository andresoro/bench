package main

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type LoadTester struct {
	endpoint string
	conns    int
	request  *http.Request
	client   *http.Client
	stats    *Stats
	ch       chan *Stats
	mu       *sync.Mutex
	dur      time.Duration
}

func NewTester(r *http.Request, conns int, dur time.Duration) *LoadTester {
	return &LoadTester{
		request: r,
		client:  &http.Client{},
		conns:   conns,
		dur:     dur,
		ch:      make(chan *Stats),
		mu:      &sync.Mutex{},
		stats:   &Stats{},
	}
}

// Run initializes the LoadTester with its # of conns for a given duration
// passes the results to the statistics channel
func (l *LoadTester) Run(ch chan *Stats) {
	var wg sync.WaitGroup

	for i := 0; i <= l.conns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for start := time.Now(); time.Since(start) < l.dur; {
				stat, err := l.test()
				if err != nil {
					ch <- &Stats{
						Endpoint: l.endpoint,
						err:      1,
					}
					continue
				}
				ch <- stat

			}
		}()
	}

	wg.Wait()
}

// Make an individual request and return the statistics
func (l *LoadTester) test() (*Stats, error) {
	var body []byte

	start := time.Now()
	resp, err := l.client.Do(l.request)
	if err != nil {
		return nil, err
	}
	end := time.Since(start)

	if resp != nil {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	s := &Stats{
		ResponseDur:  end,
		ResponseSize: float64(len(body)),
	}

	return s, nil

}
