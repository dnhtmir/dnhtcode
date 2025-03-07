package main

import (
	"flag"
	"fmt"
	"strings"

	"intruder/internal/dnhtclient"
	"intruder/internal/dnhtrequest"
	"intruder/internal/dnhtutils"
	"intruder/internal/dnhtworkers"
)

func main() {
	requestFile := flag.String("rf", "", "Path to the request file")
	rps := flag.Int("rps", 0, "Requests per second (0 for no delay)")
	jitter := flag.Float64("jitter", 0.0, "Jitter percentage (0.0 to 1.0)")
	wrk := flag.Int("wrk", 1, "Number of concurrent workers")
	values := flag.String("v", "", "Comma-separated list of values or range (e.g., 1,2,3 or 1-10)")
	flag.Parse()

	if *requestFile == "" || *values == "" {
		fmt.Println("Usage: go run main.go -rf <request_file> -values <values> [-rps <requests_per_second>] [-jitter <jitter>] [-wrk <workers>]")
		return
	}

	method, url, headers, body, err := dnhtrequest.ParseRequestFile(*requestFile)
	if err != nil {
		fmt.Println("Error parsing request file:", err)
		return
	}

	host := headers.Get("Host")
	if host == "" {
		fmt.Println("Error: no host detected")
		return
	}
	url = "https://" + host + url

	client := dnhtclient.NewClient()

	// Generate values
	var valueList []string
	if strings.Contains(*values, "-") {
		parts := strings.Split(*values, "-")
		start := dnhtutils.Atoi(parts[0])
		end := dnhtutils.Atoi(parts[1])
		for i := start; i <= end; i++ {
			valueList = append(valueList, fmt.Sprintf("%d", i))
		}
	} else {
		valueList = strings.Split(*values, ",")
	}

	jobs := make(chan string, len(valueList))
	results := make(chan string, len(valueList))

	// Start workers
	for i := range *wrk {
		worker := dnhtworkers.Worker{
			Client:  client,
			Method:  method,
			URL:     url,
			Headers: headers,
			Body:    body,
			RPS:     *rps,
			Jitter:  *jitter,
		}
		go worker.StartWorker(i, jobs, results)
	}

	// Send jobs
	for _, value := range valueList {
		jobs <- value
	}
	close(jobs)

	// Collect results
	for range len(valueList) {
		fmt.Println(<-results)
	}
}
