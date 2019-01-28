package main

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"
)

type Bench struct {
	testers map[string]*LoadTester
	ch      chan *Stats
	stats   map[string]*Stats
}

func NewBench(path string) (*Bench, error) {

	b := &Bench{
		testers: make(map[string]*LoadTester),
		ch:      make(chan *Stats),
		stats:   make(map[string]*Stats),
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
func (b *Bench) Run() {
	var wg sync.WaitGroup

	go b.handleStats()

	for _, tester := range b.testers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tester.Run(b.ch)
		}()
	}
	wg.Wait()

}

func (b *Bench) handleStats() {
	// read in statistics

	for stat := range b.ch {
		b.stats[stat.Endpoint] = stat
	}

}
