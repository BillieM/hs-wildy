package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Player contains information about all relevant categories for a player
type Player struct {
	Name   string
	Bosses map[string]Category
}

// Category contains information for a particular category for a player (e.g. chaos ele)
type Category struct {
	Name  string
	Rank  int
	Score int
}

func callAPI(playerName string) (*Player, error) {

	config := readConfig()

	var responseArr []string

	p := Player{
		Name:   playerName,
		Bosses: make(map[string]Category),
	}

	reqURL := "https://secure.runescape.com/m=hiscore_oldschool_hardcore_ironman/index_lite.ws?player=" + playerName

	resp, err := http.Get(reqURL)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	bodyS := fmt.Sprintf("%s", body)

	scanner := bufio.NewScanner(strings.NewReader(bodyS))

	for scanner.Scan() {
		lineText := scanner.Text()
		responseArr = append(responseArr, lineText)
	}

	if len(responseArr) != config.APIProperties {

		alertMsg := fmt.Sprintf("Incorrect response array length -> %s", len(responseArr))

		alertMe(alertMsg)

		return &p, errors.New("Incorrect Response Array Length")
	}

	for key, bossIndex := range config.WildernessBosses {
		bossLine := responseArr[bossIndex]
		bossLineArr := strings.Split(bossLine, ",")
		rank, err := strconv.Atoi(bossLineArr[0])
		if err != nil {
			log.Fatal(err)
		}
		score, err := strconv.Atoi(bossLineArr[1])
		if err != nil {
			log.Fatal(err)
		}
		category := Category{
			Name:  key,
			Rank:  rank,
			Score: score,
		}
		p.Bosses[key] = category
	}

	return &p, nil
}
