package main

import (
	"fmt"
	"net/http"
	"time"
)

// result represents the output of a single URL check
type result struct {
	url     string
	status  string
	latency time.Duration
	err     error
}

func main() {
	// 1. Define the targets to scan
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://caddyserver.com",
		"https://terraform.io",
		"https://thisdomainwillfail1234.com", // Intentionally broken to test error handling
	}

	// 2. Create a channel to communicate safely between goroutines
	resultsChannel := make(chan result)

	fmt.Printf("🚀 Starting concurrent status check for %d URLs...\n\n", len(urls))
	startTime := time.Now()

	// 3. Spin up a separate goroutine for each URL check
	for _, url := range urls {
		go checkStatus(url, resultsChannel)
	}

	// 4. Collect the data from the channel as soon as they finish
	for i := 0; i < len(urls); i++ {
		res := <-resultsChannel

		if res.err != nil {
			fmt.Printf("❌ [DOWN] %-30s | Error: %v\n", res.url, res.err)
		} else {
			fmt.Printf("✅ [UP]   %-30s | Status: %s | Time: %v\n", res.url, res.status, res.latency)
		}
	}

	fmt.Printf("\n✨ All checks completed in %v!\n", time.Since(startTime))
}

// checkStatus performs the HTTP GET request and passes results to the channel
func checkStatus(url string, ch chan result) {
	start := time.Now()

	// Set a reasonable timeout so it doesn't hang forever
	client := http.Client{
		Timeout: 4 * time.Second,
	}

	resp, err := client.Get(url)
	latency := time.Since(start)

	if err != nil {
		ch <- result{url: url, err: err}
		return
	}
	defer resp.Body.Close()

	ch <- result{
		url:     url,
		status:  resp.Status,
		latency: latency.Round(time.Millisecond),
		err:     nil,
	}
}
