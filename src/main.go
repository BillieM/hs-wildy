package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/ratelimit"
)

func main() {

	config := readConfig()

	db := dbConnect(config)

	runner, err := configureRunner(db)

	if err != nil {
		sendErrorAlert(fmt.Sprintf("runner configuration failed -> %v", err))
		os.Exit(0)
	}

	apiRateLimit := ratelimit.New(50, ratelimit.Per(60*time.Second))    // rate limit of 40/min
	scrapeRateLimit := ratelimit.New(10, ratelimit.Per(60*time.Second)) // rate limit of 5/min

	for i := 0; i < 3; i++ {
		go apiRunner(runner, apiRateLimit)
	}

	for i := 0; i < 2; i++ {
		go scrapeRunner(runner, scrapeRateLimit)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done // Will block here until user hits ctrl+c
}

func apiRunner(runner *Runner, rl ratelimit.Limiter) {
	for {
		rl.Take()
		runner.performApiCall()
	}
}

func scrapeRunner(runner *Runner, rl ratelimit.Limiter) {
	for {
		rl.Take()
		runner.performScrape()
	}
}
