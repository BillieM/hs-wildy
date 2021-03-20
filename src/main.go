package main

import (
	"fmt"
	"log"
	"time"
)

func createAll(db *MyDB) {
	allCategories := scrapeAll()

	var toTweet []*HighscoreLine

	for _, category := range allCategories {
		for _, page := range category.Pages {
			for _, line := range page.Lines {
				shouldTweet := db.highscoreLineCreateOrUpdate(line)
				if shouldTweet {
					toTweet = append(toTweet, line)
				}
			}
		}
	}

	for _, tweet := range toTweet {
		fmt.Println(tweet)
	}
}

func main() {

	start := time.Now()

	// fmt.Println(scrapeURL())

	// players := [2]string{"Lamui", "LarryChamp"}

	// for _, player := range players {
	// 	p, err := callAPI(player)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(p)
	// }

	db := dbConnect()
	_ = db

	createAll(db)

	// db.createPlayer("bob", true)
	// db.playerDied("bob")
	// // db.createCategory("bob", "venenatis", 2, 300)
	// db.updateCategory("bob", "venenatis", 1, 307)

	elapsed := time.Since(start)

	log.Printf("execution time %s", elapsed)
}

/*
two ways of gathering hiscores information
	scraping the hiscores
		this has to be fairly infrequent with new rate limits it seems
		but need to do some texts
	hiscores lite api
		supposedly no limit
		if kc has changed upon an api lookup, need to scrape to check player is still alive
			may require scraping surrounding pages too if player can not be found on the expected page

need some intelligent way of splitting load between the two at sensible intervals, prioritising people higher on the hiscores
	store last lookup in db,

likely going to be storing data in a db rather than pickled files

usage of structs

for each boss for each player, we store a last checked time
*/
