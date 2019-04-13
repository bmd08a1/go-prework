package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

type inputData struct {
	num_requests      int64
	concurrency_level int64
	link              string
}

type summaryInfo struct {
	hostname           string
	port               string
	path               string
	concurrency_level  int64
	requests           int64
	succeeded_requests int64
	totalTime          float64
	requestsTime       float64
}

type responseInfo struct {
	succeeded bool
	duration  float64
}

var userInput inputData
var summary summaryInfo
var start time.Time

func main() {
	initialize()
	collectInput()
	extractServerInfo()

	resultChannel := make(chan responseInfo, userInput.num_requests)
	doBenchMark(resultChannel)
	combineResult(resultChannel)

	printReport()
}

func initialize() {
	flag.Usage = func() {
		fmt.Println("Usage: mb [options] [http[s]://]hostname[:port]/path")
		fmt.Println("Options are:")
		flag.PrintDefaults()
	}
	summary = summaryInfo{"", "", "", 0, 0, 0, 0, 0}
	start = time.Now()
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
	summary.concurrency_level = userInput.concurrency_level
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

func doBenchMark(resultChannel chan responseInfo) {
	fmt.Printf("Running benchmark on %v\n", userInput.link)

	for i := int64(0); i < userInput.concurrency_level; i++ {
		go sendRequests(resultChannel)
	}
}

func sendRequests(resultChannel chan responseInfo) {
	for summary.requests < userInput.num_requests {
		requestStartAt := time.Now()

		response, _ := http.Get(userInput.link)
		summary.requests++

		resultChannel <- responseInfo{
			response.StatusCode >= 200 && response.StatusCode < 400,
			time.Now().Sub(requestStartAt).Seconds(),
		}
	}
}

func combineResult(resultChannel chan responseInfo) {
	for result := range resultChannel {
		if result.succeeded {
			summary.succeeded_requests++
		}

		summary.requestsTime += result.duration

		if summary.requests == userInput.num_requests {
			close(resultChannel)
			summary.totalTime = time.Now().Sub(start).Seconds()
		}
	}
}

func printReport() {
	fmt.Println("\nSummary:")
	fmt.Printf("Server Hostname: %v\n", summary.hostname)
	fmt.Printf("Server Port: %v\n\n", summary.port)
	fmt.Printf("Document Path: %v\n\n", summary.path)
	fmt.Printf("Document Length: %v\n\n", summary.path)
	fmt.Printf("Concurrency Level: %v\n", summary.concurrency_level)
	fmt.Printf("Requests sent: %v\n", summary.requests)
	fmt.Printf("Complete requests: %v\n", summary.succeeded_requests)
	fmt.Printf("Failed requests: %v\n", summary.requests-summary.succeeded_requests)
	fmt.Printf("Time taken for tests: %.2f (s)\n", summary.totalTime)
	fmt.Printf("Requests per second: %.2f (requests/s)\n", float64(summary.requests)/summary.totalTime)
	fmt.Printf("Time per requests: %.2f (s)\n", summary.requestsTime/float64(summary.requests))
}
