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

// APIPlayer contains information about all relevant categories for a player
type APIPlayer struct {
	Name   string
	Bosses map[string]APICategory
}

// APICategory contains information for a particular category for a player (e.g. chaos ele)
type APICategory struct {
	Name  string
	Rank  int
	Score int
}

func callAPI(playerName string) (*APIPlayer, error) {

	config := readConfig()

	var responseArr []string

	p := APIPlayer{
		Name:   playerName,
		Bosses: make(map[string]APICategory),
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

		alertMsg := fmt.Sprintf("Incorrect response array length -> %v", len(responseArr))

		sendAlert(alertMsg)

		return &p, errors.New("incorrect array response length")
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
		category := APICategory{
			Name:  key,
			Rank:  rank,
			Score: score,
		}
		p.Bosses[key] = category
	}

	return &p, nil
}
