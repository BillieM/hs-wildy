package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
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

func sendUpdateAlert(db *MyDB, changeInfo *CatChange) {

	if changeInfo.NewRank == 0 {
		writeLineToErrorLog(fmt.Sprintf("player %v has rank 0 in %v", changeInfo.PlayerName, changeInfo.CategoryName))
		return
	}

	// create a record of the change in the database

	if changeInfo.NewCategory || changeInfo.ScoreChanged {
		fmt.Println(fmt.Sprint(changeInfo))
		db.createUpdate(*changeInfo)
		writeLineToSuccessLog(fmt.Sprint(changeInfo))

		// numRecentUpdates := db.getCountRecentUpdates()

		// if numRecentUpdates < 5 {
		// 	sendTweet(changeInfo)
		// }
	}

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

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func getTwitterClient(cfg Config) *gotwi.Client {

	nci := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           cfg.AccessToken,
		OAuthTokenSecret:     cfg.AccessSecret,
	}

	client, err := gotwi.NewClient(nci)

	if err != nil {
		writeLineToErrorLog(fmt.Sprintf("error creating twitter client: %v", err))
	}

	return client
}

func sendTweet(changeInfo *CatChange) {

	config := readConfig()
	client := getTwitterClient(*config)

	resp, err := managetweet.Create(
		context.Background(),
		client,
		&types.CreateInput{
			Text: gotwi.String(changeInfo.String()),
		},
	)

	if err != nil {
		writeLineToErrorLog(fmt.Sprintf("error sending tweet: %v", err))
	}

	writeLineToSuccessLog(fmt.Sprintf("tweet sent, id: %s", gotwi.StringValue(resp.Data.ID)))

}

// func checkPreviousTweets(client *twitter.Client, changeInfo CatChange) (int64, error) {

// 	config := readConfig()

// 	uTimelineParams := twitter.UserTimelineParams{
// 		ScreenName: config.AccountName,
// 	}

// 	tweets, _, err := client.Timelines.UserTimeline(&uTimelineParams)

// 	if err != nil {
// 		return 0, err
// 	}

// 	for _, tweet := range tweets {

// 		tweetTime, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tweet.CreatedAt)
// 		if err != nil {
// 			return 0, err
// 		}
// 		minutesSinceTweet := time.Since(tweetTime).Minutes()
// 		if minutesSinceTweet < config.MinutesBetweenNewTweets {
// 			if strings.Contains(tweet.Text, changeInfo.PlayerName) && strings.Contains(tweet.Text, changeInfo.CategoryName) {
// 				return tweet.ID, nil
// 			}
// 		}
// 	}

// 	return 0, nil
// }
