package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	addr        string
	conns       int
	timeout     int
	duration    int
	header      string
	method      string
	requestData string
)

// Stats holds individual statistics from requests
type Stats struct {
	ResponseSize  int
	ResponseDur   time.Duration
	TotalRequests uint32
}

func init() {
	flag.StringVar(&addr, "a", "", "address to benchmark")
	flag.IntVar(&conns, "c", 10, "# of concurrent connections")
	flag.IntVar(&duration, "d", 5, "duration of test in seconds")
	flag.IntVar(&timeout, "t", 5, "# of seconds per request")
	flag.StringVar(&header, "h", "", "header to add to request")
	flag.StringVar(&method, "x", "GET", "http method to benchmark with")
	flag.StringVar(&requestData, "r", "", "request data to send")
	flag.Parse()

}

func main() {
	var wg sync.WaitGroup

	statsChan := make(chan *Stats, conns)
	start := time.Now()
	total := Stats{
		ResponseSize:  0,
		TotalRequests: 0,
	}

	// handle load test
	for i := 0; i <= conns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{}

			for start := time.Now(); time.Since(start) < 5*time.Second; {
				statsChan <- request(client)
			}

		}()
	}
	// channel is blocked till it is recieved, so handle all stats in this go routine
	go func() {
		for s := range statsChan {
			total.ResponseSize += s.ResponseSize
			total.ResponseDur += s.ResponseDur
			atomic.AddUint32(&total.TotalRequests, 1)
		}
	}()

	wg.Wait()
	close(statsChan)

	fmt.Printf("Average response size: %d\n", total.ResponseSize/int(total.TotalRequests))
	fmt.Printf("Requests per second %d\n", (int(total.TotalRequests) / int(time.Since(start).Seconds())))
	fmt.Printf("Total Requests sent: %d\n", int(total.TotalRequests))
	log.Printf("Test completed in %fs", time.Since(start).Seconds())

}

// request is an individual request that is sent to the server
func request(client *http.Client) *Stats {

	var buf io.Reader

	if requestData != "" {
		buf = bytes.NewBufferString(requestData)
	}

	req, err := http.NewRequest(method, addr, buf)
	if err != nil {
		log.Println("Error building new request")
		return nil
	}

	// add header
	var h []string
	if header != "" {
		h = strings.Split(header, ":")
		req.Header.Add(h[0], h[1])
	}

	// start time and make request
	var body []byte
	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request to server")
		return nil
	}
	end := time.Since(start)

	if resp != nil {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body")
		}
	}

	size := len(body)

	return &Stats{
		ResponseSize: size,
		ResponseDur:  end,
	}
}
