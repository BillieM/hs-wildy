package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// HighscoreLine stores information about a particular player on a particular hiscore page
type HighscoreLine struct {
	Rank  int
	Name  string
	Alive bool
	Score int
}

func scrapeURL() (map[string]HighscoreLine, bool) {

	highscoreLines := make(map[string]HighscoreLine)
	var isNextPage bool

	c := colly.NewCollector()

	c.OnHTML("tr.personal-hiscores__row", func(e *colly.HTMLElement) {
		playerLine := parseLine(e)
		highscoreLines[playerLine.Name] = *playerLine
	})

	c.OnHTML("a.personal-hiscores__pagination-arrow--down", func(e *colly.HTMLElement) {
		isNextPage = true
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL.String())
	})

	c.Visit("https://secure.runescape.com/m=hiscore_oldschool_hardcore_ironman/overall?category_type=1&table=20#headerHiscores")

	return highscoreLines, isNextPage

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
		Alive: alive,
	}

	return &l

}
