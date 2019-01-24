package main

import (
	"flag"
	"log"
)

var (
	file string
)

func init() {
	flag.StringVar(&file, "f", "", "config file path")
	flag.Parse()

}

func main() {
	b, err := NewBench(file)
	if err != nil {
		log.Fatal(err)
	}
	b.Run()
}
