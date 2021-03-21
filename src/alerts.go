package main

import (
	"fmt"
	"log"
	"os"
)

func writeLineToLog(msg string) {
	f, err := os.OpenFile("updates.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(msg + "\n"); err != nil {
		log.Println(err)
	}
}

func sendAlert(msg string) {
	// will send a tweet & also temporarily an email to me. (or write to file)
	fmt.Println(msg)
	writeLineToLog(msg)
}
