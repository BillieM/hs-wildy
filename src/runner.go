package main

import (
	"fmt"
)

type Runner struct {
	Categories           []string
	CurrentCategoryIndex int
	CurrentPage          int
	Database             *MyDB
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
			sendUpdateAlert(runner.Database, hsChange.Change)
		}
	}

}

func (runner *Runner) postScrapeUpdates(morePages bool) {

	if morePages && runner.CurrentPage <= 50 {
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

	reqString := fmt.Sprintf("calling api -> %s", playerName)

	fmt.Println(reqString)
	writeLineToRequestLog(reqString)

	apiData, err := callAPI(playerName)

	if err != nil {
		sendErrorAlert(err.Error())
	} else {
		runner.processAPICall(apiData)
	}
}

func (runner *Runner) processAPICall(apiData *APIPlayer) {

	if apiData.PlayerGone {
		runner.Database.playerDied(apiData.Name)
	}

	apiChanges := runner.Database.apiDataCreateOrUpdate(apiData)

	if len(apiChanges) > 0 {

		isAlive, err := scrapeIsPlayerAlive(apiData)

		if err != nil {
			sendErrorAlert(err.Error())
		} else {
			if isAlive {
				for _, catChange := range apiChanges {
					sendUpdateAlert(runner.Database, catChange)
				}
			}
		}
	}
}
