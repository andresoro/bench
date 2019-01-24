package main

import (
	"flag"
)

var (
	file string
)

func init() {
	flag.StringVar(&file, "f", "", "config file path")
	flag.Parse()

}

func main() {

}
