package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MyDB struct used to add custom methods to db
type MyDB struct {
	*gorm.DB
}

// Player struct used to generate model for players in DB
type Player struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique"`
	Alive bool
}

// Category struct used to generate model for categories in DB, each boss is a category
type Category struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	PlayerID uint
	Player   Player
	Rank     uint
	Score    uint
	Updated  time.Time `gorm:"autoUpdateTime"`
}

// Update struct used to generate model for updates in DB
type Update struct {
	ID            uint `gorm:"primaryKey"`
	CategoryName  string
	PlayerID      uint
	Player        Player
	PreviousScore uint
	NewScore      int
	PreviousRank  uint
	NewRank       int
	Updated       time.Time `gorm:"autoUpdateTime"`
}

// ScrapeData struct used to determine whether or not to tweet, and what to tweet
type CatChange struct {
	NewCategory   bool
	ScoreChanged  bool
	PlayerName    string
	CategoryName  string
	PreviousScore uint
	NewScore      int
	PreviousRank  uint
	NewRank       int
	LastUpdate    time.Time
}

type HSChange struct {
	Change      *CatChange
	PlayerAlive bool
}

func dbConnect(config *Config) *MyDB {

	var db *gorm.DB
	var err error

	HS_WILDY := os.Getenv("HSWILDY")

	if HS_WILDY == "LIVE" {
		msg := "connecting to live postgresql db"
		fmt.Println(msg)
		writeLineToOtherLog(msg)
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.DBHost, config.DBUser, config.DBPass, config.DBName)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		msg := "connecting to dev sqlite db"
		fmt.Println(msg)
		writeLineToOtherLog(msg)
		db, err = gorm.Open(sqlite.Open("../app.db"), &gorm.Config{})
	}

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Player{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&Update{})

	myDB := MyDB{
		db,
	}

	return &myDB
}

func (db *MyDB) createPlayer(name string, alive bool) *gorm.DB {
	result := db.Create(&Player{
		Name:  name,
		Alive: alive,
	})
	return result
}

func (db *MyDB) playerDied(name string) {
	db.Model(&Player{}).Where("name = ?", name).Update("alive", false)
}

func (db *MyDB) getPlayerID(playerName string) uint {
	row := db.Table("players").Where("name = ?", playerName).Select("id").Row()
	var id uint
	row.Scan(&id)
	return id
}

func (db *MyDB) createCategory(playerName string, catName string, rank uint, score uint) *gorm.DB {
	playerID := db.getPlayerID(playerName)

	result := db.Create(&Category{
		Name:     catName,
		Rank:     rank,
		Score:    score,
		PlayerID: playerID,
	})

	return result
}

func (db *MyDB) createUpdate(c CatChange) *gorm.DB {
	playerID := db.getPlayerID(c.PlayerName)

	result := db.Create(&Update{
		PlayerID:      playerID,
		CategoryName:  c.CategoryName,
		PreviousScore: c.PreviousScore,
		NewScore:      c.NewScore,
		PreviousRank:  c.PreviousRank,
		NewRank:       c.NewRank,
	})

	return result
}

// returns number of updates in the last 2 minutes
func (db *MyDB) getCountRecentUpdates() int64 {
	var count int64

	err := db.Model(&Update{}).
		Where("updated > ?", time.Now().Add(-time.Minute*2)).
		Count(&count).
		Error

	if err != nil {
		writeLineToErrorLog(fmt.Sprintf("error getting count of recent updates -> %v", err))
	}

	return count
}

func (db *MyDB) updateCategory(playerName string, catName string, playerRank uint, playerScore uint) {
	playerID := db.getPlayerID(playerName)

	db.Table("categories").Where("player_id = ? AND name = ?", playerID, catName).Updates(Category{
		Rank:  playerRank,
		Score: playerScore,
	})

}

func (db *MyDB) createOrUpdateCategory(playerName string, catName string, playerRank uint, playerScore uint) *CatChange {
	var changeData CatChange

	scoreChanged := false
	categoryExists := false

	playerID := db.getPlayerID(playerName)

	err := db.Table("categories").
		Select("count(*) > 0").
		Where("player_id = ? AND name = ?", playerID, catName).
		Find(&categoryExists).
		Error

	if err != nil {
		writeLineToErrorLog(fmt.Sprintf("error checking if category exists -> %v", err))
	}

	row := db.Table("categories").Where("name = ? AND player_id = ?", catName, playerID).Select("score", "rank", "updated").Row()

	var score uint
	var rank uint
	var updated time.Time
	row.Scan(&score, &rank, &updated)

	if !categoryExists {
		db.createCategory(
			playerName,
			catName,
			playerRank,
			playerScore,
		)
	} else {

		if score != playerScore {
			scoreChanged = true
			writeLineToOtherLog(fmt.Sprintf("db score: %v, new score: %v, player name: %s, cat name: %s", score, playerScore, playerName, catName))
		}

		db.updateCategory(playerName, catName, playerRank, playerScore)
	}

	changeData.ScoreChanged = scoreChanged
	changeData.NewCategory = !categoryExists
	changeData.PlayerName = playerName
	changeData.CategoryName = catName
	changeData.PreviousScore = score
	changeData.NewScore = int(playerScore)
	changeData.PreviousRank = rank
	changeData.NewRank = int(playerRank)
	changeData.LastUpdate = updated

	return &changeData
}

func (db *MyDB) highscoreLineCreateOrUpdate(highscoreLine *HighscoreLine) *HSChange {

	playerName := highscoreLine.Name
	playerCatRank := highscoreLine.Rank
	playerCatScore := highscoreLine.Score
	lineCatName := highscoreLine.Category

	playerIsAlive := highscoreLine.Alive

	var player Player

	playerDB := db.Where(&player, Player{
		Name: playerName,
	})

	var playerExists bool
	err := db.Model(player).
		Select("count(*) > 0").
		Where(&player, Player{
			Name: playerName,
		}).
		Find(&playerExists).
		Error

	if err != nil {
		writeLineToErrorLog(fmt.Sprintf("error checking if player exists -> %v", err))
	}

	if !playerExists {
		db.createPlayer(
			playerName,
			playerIsAlive,
		)
	} else {
		fmt.Println("player exists", playerName)
		var ID uint
		var name string
		var alive bool
		row := playerDB.Row()
		row.Scan(&ID, &name, &alive)
		if !playerIsAlive && alive {
			db.playerDied(playerName)
		}
	}

	catChange := db.createOrUpdateCategory(
		playerName,
		lineCatName,
		uint(playerCatRank),
		uint(playerCatScore),
	)

	hsChange := HSChange{
		Change:      catChange,
		PlayerAlive: playerIsAlive,
	}

	return &hsChange
}

func (db *MyDB) apiDataCreateOrUpdate(apiData *APIPlayer) []*CatChange {
	var apiChanges []*CatChange

	for _, category := range apiData.Bosses {
		if category.Score > -1 {
			catChange := db.createOrUpdateCategory(apiData.Name, category.Name, uint(category.Rank), uint(category.Score))
			if catChange.ScoreChanged || catChange.NewCategory {
				apiChanges = append(apiChanges, catChange)
			}
		}
	}

	return apiChanges
}

func (db *MyDB) getNextApiCallName() string {

	var player Player
	var name string

	qry := db.Joins("JOIN categories ON categories.player_id = players.id AND players.alive = ?", true).Order("updated").First(&player).Select("players.name").Row()
	qry.Scan(&name)

	return name
}
