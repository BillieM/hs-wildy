package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// HighscoreCategory stores information about a particular category on the hiscores, it contains highscore pages
type HighscoreCategory struct {
	Name  string
	Pages []*HighscorePage
}

// HighscorePage stores information about a particular page of a category on the hiscores, It contains highscore lines
type HighscorePage struct {
	Name  string
	Page  int
	Lines []*HighscoreLine
}

// HighscoreLine stores information about a particular player on a particular hiscore page
type HighscoreLine struct {
	Rank     int
	Name     string
	Alive    bool
	Score    int
	Category string
}

func scrapeAll() []*HighscoreCategory {
	config := readConfig()
	bosses := config.WildernessBosses

	var allCategories []*HighscoreCategory

	for bossName := range bosses {
		highscoreCat := scrapeCategory(bossName)
		allCategories = append(allCategories, highscoreCat)
	}

	return allCategories
}

func scrapeCategory(bossName string) *HighscoreCategory {

	highscoreCat := HighscoreCategory{
		Name: bossName,
	}

	pageNum := 1
	morePages := true

	for morePages {
		var highscorePage *HighscorePage
		highscorePage, morePages = scrapePage(bossName, pageNum)
		highscoreCat.Pages = append(highscoreCat.Pages, highscorePage)
		pageNum++
	}

	return &highscoreCat
}

func scrapePage(bossName string, pageNum int) (*HighscorePage, bool) {

	config := readConfig()

	tableID := config.WildernessBosses[bossName]

	url := fmt.Sprintf("https://secure.runescape.com/m=hiscore_oldschool_hardcore_ironman/overall?category_type=1&table=%v&page=%v", tableID, pageNum)

	// "https://secure.runescape.com/m=hiscore_oldschool_hardcore_ironman/overall?category_type=1&table=20#headerHiscores"

	highscorePage := HighscorePage{
		Name: bossName,
		Page: pageNum,
	}
	var isNextPage bool

	c := colly.NewCollector()

	c.OnHTML("tr.personal-hiscores__row", func(e *colly.HTMLElement) {
		playerLine := parseLine(e)
		playerLine.Category = bossName
		highscorePage.Lines = append(highscorePage.Lines, playerLine)
	})

	c.OnHTML("a.personal-hiscores__pagination-arrow--down", func(e *colly.HTMLElement) {
		isNextPage = true
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL.String())
	})

	c.Visit(url)

	return &highscorePage, isNextPage

}

func parseLine(e *colly.HTMLElement) *HighscoreLine {

	var alive bool
	var rank int
	var score int

	name := e.ChildText("td.left")

	deathImg := e.ChildAttr("img.hiscore-death", "title")

	if deathImg == "" {
		alive = true
	}

	e.ForEach("td.right", func(i int, h *colly.HTMLElement) {
		if i == 0 {
			var err error

			stringBase := h.Text
			stringTrimmed := strings.Trim(stringBase, "\n")

			rank, err = strconv.Atoi(stringTrimmed)
			if err != nil {
				log.Fatal(err)
			}
		}
		if i == 1 {
			var err error

			stringBase := h.Text
			stringNoComma := strings.Replace(stringBase, ",", "", -1)
			stringTrimmed := strings.Trim(stringNoComma, "\n")

			score, err = strconv.Atoi(stringTrimmed)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	l := HighscoreLine{
		Name:  name,
		Rank:  rank,
		Score: score,
		Alive: alive}

	return &l

}
