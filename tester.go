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
	dur      time.Duration
	ch       chan *Stats
}

func NewTester(r *http.Request, ch chan *Stats, conns int, dur time.Duration) *LoadTester {
	return &LoadTester{
		request: r,
		client:  &http.Client{},
		conns:   conns,
		dur:     dur,
		ch:      ch,
	}
}

// Run initializes the LoadTester with its # of conns for a given duration
func (l *LoadTester) Run() {
	var wg sync.WaitGroup

	for i := 0; i <= l.conns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for start := time.Now(); time.Since(start) < l.dur; {
				l.test()
			}
		}()
	}

	wg.Wait()
}

func (l *LoadTester) test() error {
	var body []byte

	start := time.Now()
	resp, err := l.client.Do(l.request)
	if err != nil {
		return err
	}
	end := time.Since(start)

	if resp != nil {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	size := len(body)

	l.ch <- &Stats{
		Endpoint:     l.endpoint,
		ResponseSize: size,
		ResponseDur:  end,
	}

	return nil

}
