package main

import (
	"fmt"
	"log"
	"os"
)

func sendUpdateAlert(msg string) {
	fmt.Println(msg)
	writeLineToSuccessLog(msg)
}

func sendErrorAlert(msg string) {
	fmt.Println(msg)
	writeLineToErrorLog(msg)
}

func writeLineToSuccessLog(msg string) {
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

func writeLineToErrorLog(msg string) {
	f, err := os.OpenFile("errors.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(msg + "\n"); err != nil {
		log.Println(err)
	}
}
