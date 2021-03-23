package main

import (
	"time"
)

// func createAll(db *MyDB) {
// 	allCategories := scrapeAll()

// 	var toTweet []string

// 	for _, category := range allCategories {
// 		for _, page := range category.Pages {
// 			for _, line := range page.Lines {
// 				newCategory, scoreChanged := db.highscoreLineCreateOrUpdate(line)
// 				if newCategory {
// 					tweet := fmt.Sprintf("%s has entered the highscores for %s. their kc is %v", line.Name, line.Category, line.Score)
// 					toTweet = append(toTweet, tweet)
// 				}

// 				if scoreChanged {
// 					tweet := fmt.Sprintf("%s's kc has changed for for %s. their kc is %v", line.Name, line.Category, line.Score)
// 					toTweet = append(toTweet, tweet)
// 				}
// 			}
// 		}
// 	}

// 	for _, tweet := range toTweet {
// 		fmt.Println(tweet)
// 	}
// }

func main() {

	config := readConfig()

	db := dbConnect()
	_ = db

	runner := configureRunner(db)

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
