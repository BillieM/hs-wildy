package main

func main() {

	// config := readConfig()

	// db := dbConnect()

	// runner, err := configureRunner(db)

	// if err != nil {
	// 	sendErrorAlert(fmt.Sprintf("runner configuration failed -> %v", err))
	// 	os.Exit(0)
	// }

	// for {

	// 	timeSinceLastScrape := time.Since(runner.LastScrapeTime).Seconds()
	// 	timeSinceLastApiCall := time.Since(runner.LastApiCallTime).Seconds()

	// 	if timeSinceLastScrape >= config.SecondsBetweenScrapes {
	// 		runner.performScrape()
	// 	}

	// 	if timeSinceLastApiCall >= config.SecondsBetweenApiCalls {
	// 		runner.performApiCall()
	// 	}

	// 	time.Sleep(250 * time.Millisecond)
	// }

	checkPreviousTweets(getTwitterClient(), "Sooo Hard RN")
}
