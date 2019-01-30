package bench

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Bench is a struct that controls the testers, channel communication
// and stat aggregation
type Bench struct {
	testers map[string]*LoadTester
	chans   map[string]chan *Stats
}

// NewBench returns a Bench tester
func NewBench(path string) (*Bench, error) {

	testers := make(map[string]*LoadTester)
	chans := make(map[string]chan *Stats)

	conf, err := fromJSON(path)
	if err != nil {
		return nil, err
	}

	for _, req := range conf.Req {
		var buf io.Reader
		addr := conf.Host + req.Endpoint

		if req.Data != "" {
			buf = bytes.NewBufferString(req.Data)
		}

		r, err := http.NewRequest(req.Method, addr, buf)
		if err != nil {
			return nil, err
		}
		// init new Tester with given request
		lt := NewTester(r, req.Connections, conf.Duration*time.Second, req.Endpoint)
		testers[req.Endpoint] = lt
		chans[req.Endpoint] = make(chan *Stats, 1)
	}

	b := &Bench{
		testers: testers,
		chans:   chans,
	}

	return b, nil
}

// Run a benchmark test with given config
// run each tester concurrently and wait for them to finish
func (b *Bench) Run() {
	var wg sync.WaitGroup

	for _, tester := range b.testers {
		wg.Add(1)

		// get channel for this endpoint
		ch := b.chans[tester.endpoint]

		fmt.Printf("Running test on %s with %d connections for %s \n", tester.endpoint, tester.conns, tester.dur.String())

		go func() {
			defer wg.Done()
			// run loadtester with the specific channel
			tester.Run(ch)
		}()

	}
	wg.Wait()

	for endp, ch := range b.chans {
		stat := <-ch
		fmt.Printf("Test completed for endpoint: %s \n", endp)
		fmt.Printf("	Total requests completed: %d \n", stat.TotalRequests)
		fmt.Printf("	Total errors: %d \n", stat.err)
		fmt.Printf("	Average response size: %f bytes\n", stat.ResponseSize)
		fmt.Printf("	Average response time: %fs \n", stat.ResponseDur.Seconds())
	}

}
