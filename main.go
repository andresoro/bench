package main

import (
	"log"
	"os"

	"github.com/andresoro/bench/bench"
)

func main() {

	arg := os.Args[1]
	b, err := bench.New(arg)
	if err != nil {
		log.Fatal("Could not open file")
	}
	b.Run()
}
