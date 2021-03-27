package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func (changeInfo CatChange) String() string {
	if changeInfo.NewCategory {
		return fmt.Sprintf("%s has entered the highscores for %s. their kc is %v.", changeInfo.PlayerName, changeInfo.CategoryName, changeInfo.NewScore)
	} else if changeInfo.ScoreChanged {
		timeSinceLastCheck := time.Since(changeInfo.LastUpdate)
		return fmt.Sprintf("%s's KC has changed for %s. their kc has increased from %v to %v. Time since this boss was last checked: %s",
			changeInfo.PlayerName, changeInfo.CategoryName, changeInfo.PreviousScore, changeInfo.NewScore, timeSinceLastCheck)
	} else {
		return "??? why did this get here"
	}
}

func checkCategoryAlert(changeInfo *CatChange) {

	if changeInfo.NewCategory {
		sendUpdateAlert(fmt.Sprint(changeInfo))
	} else if changeInfo.ScoreChanged {
		sendUpdateAlert(fmt.Sprint(changeInfo))
	}
}

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
