package main

import (
	"fmt"
	"time"
)

/*
when we scrape a page, we get back a page obj, and a bool returning if there's more pages for that

*/

type Runner struct {
	Categories           []string
	CurrentCategoryIndex int
	CurrentPage          int
	Database             *MyDB
	Scraping             bool
	CallingAPI           bool
	LastApiCallTime      time.Time
	LastScrapeTime       time.Time
}

func configureRunner(db *MyDB) (*Runner, error) {

	config := readConfig()

	var categories []string

	for category := range config.WildernessBosses {
		categories = append(categories, category)
	}

	fmt.Println(categories)

	runner := Runner{
		Categories:           categories,
		CurrentCategoryIndex: 0,
		CurrentPage:          1,
		Database:             db,
	}

	err := runner.configureConfigTableIDs()

	return &runner, err
}

func (runner *Runner) configureConfigTableIDs() error { return nil }

func (runner *Runner) performScrape() {
	runner.Scraping = true

	bossName := runner.Categories[runner.CurrentCategoryIndex]

	highscorePage, morePages := scrapePage(bossName, runner.CurrentPage)

	runner.Scraping = false
	runner.LastScrapeTime = time.Now()

	runner.postScrapeUpdates(morePages)

	runner.processPage(highscorePage)

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
		sendAlert(msg)
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

func (runner *Runner) getNextApiCallName() string {
	// get the oldest updated category where the player is alive
	// return the players name
	name := runner.Database.getNextApiCallName()
	return name
}

func (runner *Runner) performApiCall() {
	runner.CallingAPI = true
	name := runner.getNextApiCallName()
	fmt.Println("api call name -> ", name)

	runner.postAPICallUpdates()

	runner.processAPICall()
}

func (runner *Runner) postAPICallUpdates() {}

func (runner *Runner) processAPICall() {
	// going to take the api call return struct as argument
}
