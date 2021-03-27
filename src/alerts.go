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

func writeLineToLog(logName string, msg string) {
	f, err := os.OpenFile(fmt.Sprintf("../%s.log", logName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("[%s] %s\n", time.Now().Format("02/01/2006 15:04:05"), msg)); err != nil {
		log.Println(err)
	}
}

func writeLineToOtherLog(msg string) {
	writeLineToLog("other", msg)
}

func writeLineToRequestLog(msg string) {
	writeLineToLog("requests", msg)
}

func writeLineToSuccessLog(msg string) {
	writeLineToLog("updates", msg)
}

func writeLineToErrorLog(msg string) {
	writeLineToLog("errors", msg)
}
