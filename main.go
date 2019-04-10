package main

import (
	"fmt"
	"flag"
)

func main() {
	num_requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency_level := flag.Int64("c", 1, "Number of multiple requests to make at a time")
	flag.Parse()
	link := flag.Args()[0]

	fmt.Printf("number of requests: %v\n", *num_requests)
	fmt.Printf("concurrency level: %v\n", *concurrency_level)
	fmt.Printf("link: %v\n", link)
}