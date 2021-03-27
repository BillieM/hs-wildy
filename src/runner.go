package main

import (
	"fmt"
	"time"
)

type Runner struct {
	Categories           []string
	CurrentCategoryIndex int
	CurrentPage          int
	Database             *MyDB
	LastApiCallTime      time.Time
	LastScrapeTime       time.Time
}

func configureRunner(db *MyDB) (*Runner, error) {

	config := readConfig()

	var categories []string

	for category := range config.WildernessBosses {
		categories = append(categories, category)
	}

	runner := Runner{
		Categories:           categories,
		CurrentCategoryIndex: 0,
		CurrentPage:          1,
		Database:             db,
	}

	err := configureConfig()

	return &runner, err
}

func (runner *Runner) performScrape() {

	bossName := runner.Categories[runner.CurrentCategoryIndex]

	highscorePage, err := scrapePage(bossName, runner.CurrentPage)

	runner.LastScrapeTime = time.Now()

	if err != nil {
		sendErrorAlert(fmt.Sprintf("scrape failed -> %v", err))
	} else {
		runner.postScrapeUpdates(highscorePage.MorePages)
		runner.processPage(highscorePage)
	}

}

func (runner *Runner) processPage(highscorePage *HighscorePage) {

	for _, line := range highscorePage.Lines {
		hsChange := runner.Database.highscoreLineCreateOrUpdate(line)

		if hsChange.PlayerAlive {
			checkCategoryAlert(hsChange.Change)
		}
	}

}

func (runner *Runner) postScrapeUpdates(morePages bool) {

	if morePages {
		runner.CurrentPage++
	} else {
		runner.CurrentPage = 1
		if runner.CurrentCategoryIndex >= len(runner.Categories)-1 {
			runner.CurrentCategoryIndex = 0
		} else {
			runner.CurrentCategoryIndex++
		}
	}
}

func (runner *Runner) performApiCall() {
	playerName := runner.Database.getNextApiCallName()

	fmt.Println("api call name -> ", playerName)

	apiData, err := callAPI(playerName)

	runner.LastApiCallTime = time.Now()

	if err != nil {
		sendErrorAlert(err.Error())
	} else {
		runner.processAPICall(apiData)
	}
}

func (runner *Runner) processAPICall(apiData *APIPlayer) {

	apiChanges := runner.Database.apiDataCreateOrUpdate(apiData)

	if len(apiChanges) > 0 {

		isAlive, err := scrapeIsPlayerAlive(apiData)

		if err != nil {
			sendErrorAlert(err.Error())
		} else {
			runner.LastScrapeTime = time.Now()

			if isAlive {
				for _, catChange := range apiChanges {
					checkCategoryAlert(catChange)
				}
			}
		}
	}
}
