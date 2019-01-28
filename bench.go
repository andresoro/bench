package main

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
	ch      chan *Stats
	stats   *Stats
}

// NewBench returns a Bench tester
func NewBench(path string) (*Bench, error) {

	b := &Bench{
		testers: make(map[string]*LoadTester),
		ch:      make(chan *Stats),
		stats:   &Stats{},
	}

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
		lt := NewTester(r, req.Connections, 5*time.Second)
		b.testers[req.Endpoint] = lt

	}

	return b, nil
}

// Run a benchmark test with given config
// run each tester concurrently and wait for them to finish
func (b *Bench) Run() {
	var wg sync.WaitGroup

	go b.handleStats()

	for _, tester := range b.testers {
		wg.Add(1)
		go func(ch chan *Stats) {
			defer wg.Done()
			tester.Run(ch)
		}(b.ch)

	}

	wg.Wait()
	close(b.ch)

	fmt.Printf("Total Requests: %d \n", b.stats.TotalRequests)
	fmt.Printf("Average Response Size: %f \n", b.stats.ResponseSize)
	fmt.Printf("Average Request Time: %d \n", b.stats.ResponseDur)
	fmt.Printf("Total Errors: %d \n", b.stats.err)
}

func (b *Bench) handleStats() {

	for stat := range b.ch {
		b.stats.err += stat.err
		b.stats.ResponseDur += stat.ResponseDur
		b.stats.ResponseSize += stat.ResponseSize
		b.stats.TotalRequests++
	}

	// take averages

	b.stats.ResponseSize = b.stats.ResponseSize / float64(b.stats.TotalRequests)
	b.stats.ResponseDur = b.stats.ResponseDur / time.Duration(b.stats.TotalRequests)

}
