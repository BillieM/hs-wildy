package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func (changeInfo CatChange) String() string {
	if changeInfo.NewCategory {
		return fmt.Sprintf("%s has entered the highscores for %s (rank %v). Their kc is %v.", changeInfo.PlayerName, changeInfo.CategoryName, changeInfo.NewRank, changeInfo.NewScore)
	} else if changeInfo.ScoreChanged {
		timeSinceLastCheck := fmtDuration(time.Since(changeInfo.LastUpdate))
		return fmt.Sprintf("%s's KC has changed for %s (rank %v -> %v). Their kc has increased from %v to %v. Time since this boss was last checked: %s",
			changeInfo.PlayerName, changeInfo.CategoryName, changeInfo.PreviousRank, changeInfo.NewRank, changeInfo.PreviousScore, changeInfo.NewScore, timeSinceLastCheck)
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
	sendTweet(msg)
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

func fmtDuration(d time.Duration) string {
	numSeconds := int(d.Seconds())
	s := numSeconds % 60
	m := int(math.Floor(float64(numSeconds)/60)) % 60
	h := math.Floor(float64(numSeconds) / 3600)

	if h > 0 {
		return fmt.Sprintf("%vh:%vm:%vs", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%vm:%vs", m, s)
	} else {
		return fmt.Sprintf("%vs", s)
	}
}

func getTwitterClient() *twitter.Client {

	config := readConfig()

	oAuthCfg := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	oAuthTkn := oauth1.NewToken(config.AccessToken, config.AccessSecret)
	httpClient := oAuthCfg.Client(oauth1.NoContext, oAuthTkn)
	client := twitter.NewClient(httpClient)

	return client
}

func sendTweet(msg string) {

	client := getTwitterClient()

	_, _, err := client.Statuses.Update(msg, nil)

	if err != nil {
		sendErrorAlert(err.Error())
	}
}
