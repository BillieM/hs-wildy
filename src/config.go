package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

// Config struct stores configuration properties for application
type Config struct {
	APIProperties           int
	ScrapeProperties        int
	WildernessBosses        map[string]int
	WildernessBossesArr     []string
	SecondsBetweenScrapes   float64
	SecondsBetweenApiCalls  float64
	MinutesBetweenNewTweets float64
	NumSkills               int
	AccessToken             string
	AccessSecret            string
	AccountName             string

	DBHost string
	DBName string
	DBUser string
	DBPass string
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
	config.APIProperties = highscoreCatsInfo.NumHighscoreCategories + 2
	config.ScrapeProperties = highscoreCatsInfo.NumHighscoreCategories + 1

	config.AccountName = "HcWildyTest"
	if os.Getenv("PRODUCTION") == "TRUE" {
		config.AccountName = "HcWildy"
	}

	if err != nil {
		return err
	}

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

	err = getSecrets(&config)

	if err != nil {
		log.Fatal(err)
	}

	return &config
}

func writeConfig(config *Config) error {

	config.AccessToken = ""
	config.AccessSecret = ""
	config.DBHost = ""
	config.DBName = ""
	config.DBUser = ""
	config.DBPass = ""

	jsonData, err := json.Marshal(config)

	if err != nil {
		return err
	}

	err = os.WriteFile("../config.json", jsonData, 0644)

	if err != nil {
		return err
	}

	return nil

}

func getSecrets(config *Config) error {
	accessToken, accessTokenExists := os.LookupEnv("HCWILDY_ACCESS_TOKEN")
	accessSecret, accessSecretExists := os.LookupEnv("HCWILDY_ACCESS_SECRET")

	if !accessTokenExists || !accessSecretExists {
		return errors.New("twitter access token or secret missing")
	}

	_, consumerExists := os.LookupEnv("GOTWI_API_KEY")
	_, consumerSecretExists := os.LookupEnv("GOTWI_API_KEY_SECRET")

	if !consumerExists || !consumerSecretExists {
		return errors.New("twitter consumer key or secret missing")
	}

	config.AccessToken = accessToken
	config.AccessSecret = accessSecret

	dbHost, i1 := os.LookupEnv("HCWILDY_DB_HOST")
	dbName, i2 := os.LookupEnv("HCWILDY_DB_NAME")
	dbUser, i3 := os.LookupEnv("HCWILDY_DB_USER")
	dbPass, i4 := os.LookupEnv("HCWILDY_DB_PASS")

	if os.Getenv("HSWILDY") == "LIVE" {
		if i1 && i2 && i3 && i4 {
			// all env vars exist
			config.DBHost = dbHost
			config.DBName = dbName
			config.DBUser = dbUser
			config.DBPass = dbPass
		} else {
			return errors.New("database creds missing on production env")
		}
	}

	return nil
}
