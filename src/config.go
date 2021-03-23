package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config struct stores configuration properties for application
type Config struct {
	APIProperties          int
	WildernessBosses       map[string]int
	SecondsBetweenScrapes  float64
	SecondsBetweenApiCalls float64
}

func readConfig() *Config {
	file, err := os.Open("../config.json")

	if err != nil {
		log.Fatal(err)
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
