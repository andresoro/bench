package main

import (
	"flag"
	"log"

	"github.com/andresoro/bench/bench"
)

var (
	file string
)

func init() {
	flag.StringVar(&file, "f", "", "config file path")
	flag.Parse()

}

func main() {
	b, err := bench.New(file)
	if err != nil {
		log.Fatal(err)
	}
	b.Run()
}
