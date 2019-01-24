package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	file string
)

func init() {
	flag.StringVar(&file, "f", "", "config file path")
	flag.Parse()

}

func main() {
	var wg sync.WaitGroup
	var statChannels map[string]chan *Stats

	if file != "" {
		conf, err := fromJSON(file)
		if err != nil {
			log.Fatal(err)
		}
		// initialize map endpoint:channel
		for _, req := range conf.req {
			statChannels[req.Endpoint] = make(chan *Stats)
		}

	}

	total := Stats{
		ResponseSize:  0,
		TotalRequests: 0,
	}

	start := time.Now()

	// handle load test for each endpoint
	for i := 0; i <= conns; i++ {
		wg.Add(1)
		go func() {
			c := &http.Client{}

		}()
	}
	// channel is blocked till it is recieved, so handle all stats in this go routine
	// divide based on request
	go func() {

	}()

	wg.Wait()
	for _, ch := range statChannels {
		close(ch)
	}

	fmt.Printf("Average response size: %d\n", total.ResponseSize/int(total.TotalRequests))
	fmt.Printf("Requests per second %d\n", (int(total.TotalRequests) / int(time.Since(start).Seconds())))
	fmt.Printf("Total Requests sent: %d\n", int(total.TotalRequests))
	log.Printf("Test completed in %fs", time.Since(start).Seconds())

}

// request is an individual request that is sent to the server
func request(client *http.Client, req *http.Request) *Stats {

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
