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

	var toAlert []string

	for _, line := range highscorePage.Lines {
		changeData := runner.Database.highscoreLineCreateOrUpdate(line)

		if changeData.PlayerAlive {
			if changeData.NewCategory {
				msg := fmt.Sprintf("%s has entered the highscores for %s. their kc is %v.", line.Name, line.Category, line.Score)
				toAlert = append(toAlert, msg)
			} else if changeData.ScoreChanged {
				timeSinceLastCheck := time.Since(changeData.LastUpdate)
				msg := fmt.Sprintf("%s's KC has changed for %s. their kc has increased from %v to %v. Time since this player was last checked: %s", line.Name, line.Category, changeData.PreviousScore, line.Score, timeSinceLastCheck)
				toAlert = append(toAlert, msg)
			}
		}
	}

	for _, msg := range toAlert {
		sendUpdateAlert(msg)
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

	fmt.Println(apiData)

	runner.LastApiCallTime = time.Now()

	if err != nil {
		sendErrorAlert(err.Error())
	} else {
		runner.processAPICall(apiData)
	}
}

func (runner *Runner) processAPICall(apiData *APIPlayer) {
	// going to take the api call return struct as argument

	// if there is a change, we need to perform a scrape to check if the player is alive
	// 	check the hs page based on players rank from api call

	var bossName string
	var rank int

	for k, v := range apiData.Bosses {
		bossName = k
		rank = v.Rank
		if rank > 0 {
			break
		}
	}

	isAlive, err := scrapeIsPlayerAlive(apiData.Name, bossName, rank)

	if err != nil {
		sendErrorAlert(err.Error())
	}

	if isAlive {
		fmt.Println("alive xx")
	}

}
