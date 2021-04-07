package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config struct stores configuration properties for application
type Config struct {
	APIProperties          int
	ScrapeProperties       int
	WildernessBosses       map[string]int
	WildernessBossesArr    []string
	SecondsBetweenScrapes  float64
	SecondsBetweenApiCalls float64
	NumSkills              int
}

func configureConfig() error {

	config := readConfig()

	var categories []string

	for k := range config.WildernessBosses {
		categories = append(categories, k)
	}

	highscoreCatsInfo, err := scrapeCategoriesInfo(categories)

	if err != nil {
		return err
	}

	config.WildernessBosses = highscoreCatsInfo.CategoryIDs
	config.APIProperties = highscoreCatsInfo.NumHighscoreCategories + 1
	config.ScrapeProperties = highscoreCatsInfo.NumHighscoreCategories

	err = writeConfig(config)

	if err != nil {
		return err
	}

	return nil

}

func readConfig() *Config {
	file, err := os.Open("../config.json")

	if err != nil {
		sendErrorAlert(err.Error())
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	config := Config{}

	err = decoder.Decode(&config)

	if err != nil {
		log.Fatal(err)
	}

	return &config
}

func writeConfig(config *Config) error {

	jsonData, err := json.Marshal(config)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("../config.json", jsonData, 0644)

	if err != nil {
		return err
	}

	return nil

}
