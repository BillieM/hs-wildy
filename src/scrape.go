package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// HighscoreCategory stores information about a particular category on the hiscores, it contains highscore pages
type HighscoreCategory struct {
	Name  string
	Pages []*HighscorePage
}

// HighscorePage stores information about a particular page of a category on the hiscores, It contains highscore lines
type HighscorePage struct {
	Name      string
	Page      int
	Lines     []*HighscoreLine
	MorePages bool
}

// HighscoreLine stores information about a particular player on a particular hiscore page
type HighscoreLine struct {
	Rank     int
	Name     string
	Alive    bool
	Score    int
	Category string
}

// HighscoreCategoriesInfo stores information about categories in the hiscore, to be used to update the config file
type HighscoreCategoriesInfo struct {
	CategoryIDs            map[string]int
	NumHighscoreCategories int
}

func scrapeCategoriesInfo(categories []string) (*HighscoreCategoriesInfo, error) {

	config := readConfig()

	highscoreCatsInfo := HighscoreCategoriesInfo{}

	categoriesMap := make(map[string]int)

	url := "https://secure.runescape.com/m=hiscore_oldschool/overall"

	c := colly.NewCollector()

	var err error = nil

	c.OnHTML("div#contentCategory", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, h *colly.HTMLElement) {
			for _, category := range config.WildernessBossesArr {
				if strings.Contains(category, h.Text) {
					href := h.Attr("href")
					re := regexp.MustCompile(`table=\s*(\d+)`)
					match := re.FindStringSubmatch(href)
					if match != nil {
						var id int
						id, err = strconv.Atoi(match[1])
						categoriesMap[h.Text] = id
					} else {
						err = errors.New("unable to parse table id")
					}
				}
			}
		})
		highscoreCatsInfo.CategoryIDs = categoriesMap
		highscoreCatsInfo.NumHighscoreCategories = numChildAttrs(e)
	})

	c.OnRequest(func(r *colly.Request) {
		reqString := "scraping -> table ids"
		fmt.Println(reqString)
		writeLineToRequestLog(reqString)
	})

	c.OnError(func(_ *colly.Response, reqErr error) {
		err = reqErr
	})

	c.Visit(url)

	return &highscoreCatsInfo, err
}

func scrapePage(bossName string, pageNum int) (*HighscorePage, error) {

	config := readConfig()

	tableID := config.WildernessBosses[bossName]

	url := fmt.Sprintf("https://secure.runescape.com/m=hiscore_oldschool_hardcore_ironman/overall?category_type=1&table=%v&page=%v", tableID, pageNum)

	var err error

	highscorePage := HighscorePage{
		Name:      bossName,
		Page:      pageNum,
		MorePages: false,
	}

	c := colly.NewCollector()

	c.OnHTML("div#contentCategory", func(e *colly.HTMLElement) {
		numCats := numChildAttrs(e)
		if numCats != config.ScrapeProperties {
			err = errors.New("incorrect number of scrape categories")
			configureConfig()
		}
	})

	c.OnHTML("tr.personal-hiscores__row", func(e *colly.HTMLElement) {
		var playerLine *HighscoreLine
		playerLine, err = parseLine(e)
		playerLine.Category = bossName
		highscorePage.Lines = append(highscorePage.Lines, playerLine)
	})

	c.OnHTML("a.personal-hiscores__pagination-arrow--down", func(e *colly.HTMLElement) {
		highscorePage.MorePages = true
	})

	c.OnRequest(func(r *colly.Request) {
		reqString := fmt.Sprintf("scraping -> %s page %v", bossName, pageNum)
		fmt.Println(reqString)
		writeLineToRequestLog(reqString)
	})

	c.OnError(func(_ *colly.Response, reqErr error) {
		err = reqErr
	})

	c.Visit(url)

	return &highscorePage, err

}

func parseLine(e *colly.HTMLElement) (*HighscoreLine, error) {

	var alive bool
	var rank int
	var score int

	l := HighscoreLine{}

	var err error

	name := e.ChildText("td.left")

	deathImg := e.ChildAttr("img.hiscore-death", "title")

	if deathImg == "" {
		alive = true
	}

	e.ForEach("td.right", func(i int, h *colly.HTMLElement) {
		if i == 0 {
			stringBase := h.Text
			stringTrimmed := strings.Trim(stringBase, "\n")

			rank, err = strconv.Atoi(stringTrimmed)
		}
		if i == 1 {
			stringBase := h.Text
			stringNoComma := strings.Replace(stringBase, ",", "", -1)
			stringTrimmed := strings.Trim(stringNoComma, "\n")

			score, err = strconv.Atoi(stringTrimmed)
		}
	})

	l.Name = name
	l.Rank = rank
	l.Score = score
	l.Alive = alive

	return &l, err

}

func scrapeIsPlayerAlive(apiData *APIPlayer) (bool, error) {
	var err error
	var alive bool

	var catRank int
	var catName string

	playerName := apiData.Name

	for _, category := range apiData.Bosses {
		if category.Score > -1 {
			catRank = category.Rank
			catName = category.Name
			break
		}
	}

	pageToScrape := int(math.Ceil(float64(catRank) / 25))

	highscorePage, err := scrapePage(catName, pageToScrape)

	if err != nil {
		return alive, err
	}

	numCycles := 0

	for _, highscoreLine := range highscorePage.Lines {
		if highscoreLine.Name == playerName {
			alive = highscoreLine.Alive
			break
		} else {
			numCycles++
		}
	}

	if numCycles == 25 {
		return alive, errors.New("player not on page")
	}

	return alive, nil
}

func numChildAttrs(e *colly.HTMLElement) int {
	numChildren := 0
	e.ForEach("a", func(_ int, h *colly.HTMLElement) {
		numChildren++
	})
	return numChildren
}
