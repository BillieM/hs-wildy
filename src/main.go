package main

import (
	"time"
)

func main() {

	config := readConfig()

	db := dbConnect()
	_ = db

	runner, err := configureRunner(db)

	if err != nil {
		// exit here
	}

	for {
		//
		timeSinceLastScrape := time.Since(runner.LastScrapeTime).Seconds()
		timeSinceLastApiCall := time.Since(runner.LastApiCallTime).Seconds()

		if !runner.Scraping && timeSinceLastScrape >= config.SecondsBetweenScrapes {
			runner.performScrape()
		}

		if !runner.CallingAPI && timeSinceLastApiCall >= config.SecondsBetweenApiCalls {
			runner.performApiCall()
		}

		time.Sleep(250 * time.Millisecond)
	}
}
