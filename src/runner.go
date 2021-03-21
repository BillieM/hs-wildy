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
	LastScrapeTime       time.Time
}

func configureRunner(db *MyDB) *Runner {
	config := readConfig()

	var categories []string

	for k := range config.WildernessBosses {
		categories = append(categories, k)
	}

	runner := Runner{
		Categories:           categories,
		CurrentCategoryIndex: 0,
		CurrentPage:          1,
		Database:             db,
	}

	return &runner
}

func (runner *Runner) performScrape() {
	runner.Scraping = true

	bossName := runner.Categories[runner.CurrentCategoryIndex]

	highscorePage, morePages := scrapePage(bossName, runner.CurrentPage)

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
	runner.Scraping = false
	runner.LastScrapeTime = time.Now()

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
