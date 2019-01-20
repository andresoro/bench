package main

import (
	"flag"
	"runtime"
)

var (
	addr     string
	conns    int
	threads  int
	duration string
	header   string
)

func init() {
	flag.StringVar(&addr, "a", "", "address to benchmark")
	flag.IntVar(&conns, "c", 20, "# of concurrent connections")
	flag.IntVar(&threads, "t", runtime.GOMAXPROCS(-1), "# of threads to use")
	flag.StringVar(&duration, "d", "5", "# of seconds per request")
	flag.StringVar(&header, "h", "", "header to add to request")
	flag.Parse()

}

func main() {

}
