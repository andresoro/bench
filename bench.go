package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	addr        string
	conns       int
	timeout     int
	header      string
	method      string
	requestData string

	client *http.Client
)

func init() {
	flag.StringVar(&addr, "a", "", "address to benchmark")
	flag.IntVar(&conns, "c", 10, "# of concurrent connections")
	flag.IntVar(&timeout, "t", 5, "# of seconds per request")
	flag.StringVar(&header, "h", "", "header to add to request")
	flag.StringVar(&method, "x", "GET", "http method to benchmark with")
	flag.StringVar(&requestData, "r", "", "request data to send")
	flag.Parse()

	client = &http.Client{}

}

func main() {
	var wg sync.WaitGroup

	stats := make(chan Stats, conns)

	wg.Add(conns)

	for i := 0; i <= conns; i++ {
		go func() {
			defer wg.Done()
			request(stats)
		}()
	}

	wg.Wait()

}

// request is meant to be called many times concurrently
// results are sent to the Stats channel
func request(ch chan Stats) {

	var buf io.Reader

	if requestData != "" {
		buf = bytes.NewBufferString(requestData)
	}

	req, err := http.NewRequest(method, addr, buf)
	if err != nil {
		log.Println("Error building new request")
		return
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
		log.Println(err)
		log.Println("Error making request to server")
		return
	}

	if resp != nil {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body")
		}
	}

	// read request and calculate size

	size := len(body)
	end := time.Since(start)

	stats := Stats{
		ResponseSize: size,
		ResponseDur:  end,
	}

	log.Printf("Request made with size: %d and duration: %s", size, end)

	ch <- stats
	return

}
