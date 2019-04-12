package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

type summaryInfo struct {
	hostname          string
	port              string
	path              string
	concurrency_level int64
	requests          int64
	success_requests  int64
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: mb [options] [http[s]://]hostname[:port]/path")
		fmt.Println("Options are:")
		flag.PrintDefaults()
	}
	num_requests, concurrency_level, link := collectInput()
	hostname, port, path := extractServerInfo(link)

	summary := summaryInfo{
		hostname,
		port,
		path,
		concurrency_level,
		0,
		0,
	}

	fmt.Printf("Running benchmark on %v\n", link)

	for summary.requests < num_requests {
		summary.requests++
		response, _ := http.Get(link)
		if response.StatusCode >= 200 && response.StatusCode < 400 {
			summary.success_requests++
		}
	}

	fmt.Println("\nSummary:")
	fmt.Printf("Server Hostname: %v\n", summary.hostname)
	fmt.Printf("Server Port: %v\n", summary.port)
	fmt.Printf("Document Path: %v\n", summary.path)
	fmt.Printf("Concurrency Level: %v\n", summary.concurrency_level)
}

func extractServerInfo(link string) (hostname string, port string, path string) {
	url_regex := regexp.MustCompile(`(?P<protocol>http|https):\/\/(?P<hostname>[\w\.\-_]+):?(?P<port>\d*)(?P<path>[\w\W]*)`)
	subMatches := url_regex.FindStringSubmatch(link)

	serverInfo := make(map[string]string)

	for i, name := range url_regex.SubexpNames() {
		if i != 0 && name != "" {
			serverInfo[name] = subMatches[i]
		}
	}

	hostname = serverInfo["hostname"]
	if serverInfo["port"] != "" {
		port = serverInfo["port"]
	} else {
		if serverInfo["protocol"] == "http" {
			port = "80"
		} else {
			port = "443"
		}
	}
	path = serverInfo["path"]

	return
}

func collectInput() (int64, int64, string) {
	num_requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency_level := flag.Int64("c", 1, "Number of multiple requests to make at a time")
	flag.Parse()

	if *num_requests < 0 {
		fmt.Printf("invalid value %v for flag -n: must be positive\n", *num_requests)
		flag.Usage()
		os.Exit(1)
	}
	if *concurrency_level < 0 {
		fmt.Printf("invalid value %v for flag -c: must be positive\n", *concurrency_level)
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() < 1 {
		fmt.Println("wrong number of arguments")
		flag.Usage()
		os.Exit(1)
	}

	link := flag.Args()[0]
	url_regex, _ := regexp.Compile(`(http|https):\/\/([\w\.\-_]+):?(\d?)[\w\W]*`)
	if url_regex.FindString(link) == "" {
		fmt.Printf("invalid value '%v': must be a valid URI\n", link)
		flag.Usage()
		os.Exit(1)
	}

	return *num_requests, *concurrency_level, link
}
