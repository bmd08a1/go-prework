package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

type inputData struct {
	num_requests      int64
	concurrency_level int64
	link              string
}

type summaryInfo struct {
	hostname          string
	port              string
	path              string
	concurrency_level int64
	requests          int64
	success_requests  int64
}

var userInput inputData
var summary summaryInfo

func main() {
	initialize()
	collectInput()
	extractServerInfo()

	fmt.Printf("Running benchmark on %v\n", userInput.link)

	sendRequests(userInput.num_requests)

	printReport()
}

func initialize() {
	flag.Usage = func() {
		fmt.Println("Usage: mb [options] [http[s]://]hostname[:port]/path")
		fmt.Println("Options are:")
		flag.PrintDefaults()
	}
	summary = summaryInfo{"", "", "", 0, 0, 0}
}

func collectInput() {
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
	url_regex, _ := regexp.Compile(`(http|https):\/\/([\w\.\-_]+):?(\d*)[\w\W]*`)
	if url_regex.FindString(link) == "" {
		fmt.Printf("invalid value '%v': must be a valid URI\n", link)
		flag.Usage()
		os.Exit(1)
	}

	userInput = inputData{
		*num_requests,
		*concurrency_level,
		link,
	}
}

func extractServerInfo() {
	regex_string := `(?P<protocol>http|https):\/\/(?P<hostname>[\w\.\-_]+):?(?P<port>\d*)(?P<path>[\w\W]*)`
	url_regex := regexp.MustCompile(regex_string)
	subMatches := url_regex.FindStringSubmatch(userInput.link)

	serverInfo := make(map[string]string)

	for i, name := range url_regex.SubexpNames() {
		if i != 0 && name != "" {
			serverInfo[name] = subMatches[i]
		}
	}

	summary.hostname = serverInfo["hostname"]
	if serverInfo["port"] != "" {
		summary.port = serverInfo["port"]
	} else {
		if serverInfo["protocol"] == "http" {
			summary.port = "80"
		} else {
			summary.port = "443"
		}
	}
	summary.path = serverInfo["path"]
}

func sendRequests(num_requests int64) {
	for summary.requests < num_requests {
		summary.requests++
		response, _ := http.Get(userInput.link)
		if response.StatusCode >= 200 && response.StatusCode < 400 {
			summary.success_requests++
		}
	}
}

func printReport() {
	fmt.Println("\nSummary:")
	fmt.Printf("Server Hostname: %v\n", summary.hostname)
	fmt.Printf("Server Port: %v\n", summary.port)
	fmt.Printf("Document Path: %v\n", summary.path)
	fmt.Printf("Concurrency Level: %v\n", summary.concurrency_level)
}
