package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	config := readConfig()

	db := dbConnect(config)

	runner, err := configureRunner(db)

	if err != nil {
		sendErrorAlert(fmt.Sprintf("runner configuration failed -> %v", err))
		os.Exit(0)
	}

	for {

		timeSinceLastScrape := time.Since(runner.LastScrapeTime).Seconds()
		timeSinceLastApiCall := time.Since(runner.LastApiCallTime).Seconds()

		if timeSinceLastScrape >= config.SecondsBetweenScrapes {
			runner.performScrape()
		}

		if timeSinceLastApiCall >= config.SecondsBetweenApiCalls {
			runner.performApiCall()
		}

		time.Sleep(250 * time.Millisecond)
	}

	// prevTweet, err := checkPreviousTweets(getTwitterClient(), CatChange{
	// 	PlayerName:   "ydanus",
	// 	CategoryName: "Callisto",
	// })
	// if err != nil {
	// 	sendErrorAlert(err.Error())
	// }
	// fmt.Println(prevTweet)
	// client := getTwitterClient()
	// _, _, err = client.Statuses.Update("@HcWildy test", &twitter.StatusUpdateParams{
	// 	InReplyToStatusID: prevTweet,
	// })
	// if err != nil {
	// 	sendErrorAlert(err.Error())
	// }
}
