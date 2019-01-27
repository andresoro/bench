package main

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type Bench struct {
	testers map[string]*LoadTester
	ch      chan *Stats
	conf    config
}

func NewBench(path string) (*Bench, error) {

	b := &Bench{
		testers: make(map[string]*LoadTester),
		ch:      make(chan *Stats),
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

func (b *Bench) Run() {
	for _, tester := range b.testers {
		tester.Run(b.ch)
	}

}

func (b *Bench) handleStats() {
	// read in statistics

}
