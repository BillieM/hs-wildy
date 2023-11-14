package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// APIPlayer contains information about all relevant categories for a player
type APIPlayer struct {
	Name       string
	PlayerGone bool
	Bosses     map[string]APICategory
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

	spacesPlayerName := strings.ReplaceAll(playerName, string(rune(160)), " ")
	encodedPlayerName := url.QueryEscape(spacesPlayerName)
	reqURL := "https://secure.runescape.com/m=hiscore_oldschool_hardcore_ironman/index_lite.ws?player=" + encodedPlayerName

	resp, err := http.Get(reqURL)

	if err != nil {
		return &p, err
	}

	if resp.StatusCode == 404 {
		p.PlayerGone = true
		return &p, nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return &p, err
	}

	bodyS := string(body)

	scanner := bufio.NewScanner(strings.NewReader(bodyS))

	for scanner.Scan() {
		lineText := scanner.Text()
		responseArr = append(responseArr, lineText)
	}

	if len(responseArr) != config.APIProperties {
		configureConfig()
		return &p, errors.New("incorrect number of api categories")
	}

	for key, bossIndex := range config.WildernessBosses {
		bossLine := responseArr[bossIndex+config.NumSkills+1]
		bossLineArr := strings.Split(bossLine, ",")
		rank, err := strconv.Atoi(bossLineArr[0])
		if err != nil {
			return &p, err
		}

		if rank == 0 {
			writeLineToErrorLog(fmt.Sprintf("player %v has rank 0 in %v", playerName, key))
			continue
		}

		score, err := strconv.Atoi(bossLineArr[1])
		if err != nil {
			return &p, err
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
