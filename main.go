package main

import (
	"fmt"
	"flag"
	"os"
	"regexp"
)

func main() {
	num_requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency_level := flag.Int64("c", 1, "Number of multiple requests to make at a time")
	flag.Parse()

	if (*num_requests < 0) {
		fmt.Printf("invalid value %v for flag -n: must be positive\n", *num_requests)
		flag.PrintDefaults()
		os.Exit(1)
	}
	if (*concurrency_level < 0) {
		fmt.Printf("invalid value %v for flag -c: must be positive\n", *concurrency_level)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if (flag.NArg() < 1) {
		fmt.Println("wrong number of arguments")
		flag.PrintDefaults()
		os.Exit(1)
	}
	link := flag.Args()[0]
	url_regex := regexp.Compile(`(http|https):\/\/[\w\-_]+(?:(?:\.[\w\-_]+)+)([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	if (url_regex.FindString(link) == "") {
		fmt.Printf("invalid value '%v': must be a valid URI\n", link)
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("number of requests: %v\n", *num_requests)
	fmt.Printf("concurrency level: %v\n", *concurrency_level)
	fmt.Printf("link: %v\n", link)
}