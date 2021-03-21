package main

import (
	"errors"
	"log"
	"time"

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

func dbConnect() *MyDB {
	db, err := gorm.Open(sqlite.Open("../app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Player{})
	db.AutoMigrate(&Category{})

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

func (db *MyDB) updateCategory(playerName string, catName string, playerRank uint, playerScore uint) {
	playerID := db.getPlayerID(playerName)

	db.Table("categories").Where("player_id = ? AND name = ?", playerID, catName).Updates(Category{
		Rank:  playerRank,
		Score: playerScore,
	})

}

func (db *MyDB) createOrUpdateCategory(playerName string, catName string, playerRank uint, playerScore uint) bool {

	var category Category

	newCategory := false
	scoreChanged := false

	playerID := db.getPlayerID(playerName)

	catDB := db.Table("categories").Where("player_id = ? AND name = ?", playerID, catName).First(&category)

	newCategory = errors.Is(catDB.Error, gorm.ErrRecordNotFound)

	if newCategory {
		db.createCategory(
			playerName,
			catName,
			playerRank,
			playerScore,
		)
	} else {

		row := db.Table("categories").Where("name = ? AND player_id = ?", catName, playerID).Select("score").Row()
		var score uint
		row.Scan(&score)

		if score != playerScore {
			scoreChanged = true
		}
		db.updateCategory(playerName, catName, playerRank, playerScore)
	}

	return newCategory || scoreChanged
}

func (db *MyDB) highscoreLineCreateOrUpdate(highscoreLine *HighscoreLine) bool {

	newPlayer := false
	scoreChanged := false

	playerName := highscoreLine.Name
	playerCatRank := highscoreLine.Rank
	playerCatScore := highscoreLine.Score
	lineCatName := highscoreLine.Category

	playerIsAlive := highscoreLine.Alive

	var player Player

	playerDB := db.First(&player, Player{
		Name: playerName,
	})

	newPlayer = errors.Is(playerDB.Error, gorm.ErrRecordNotFound)

	if newPlayer {
		db.createPlayer(
			playerName,
			playerIsAlive,
		)
	} else {
		var ID uint
		var name string
		var alive bool
		row := playerDB.Row()
		row.Scan(&ID, &name, &alive)
		if !playerIsAlive && alive {
			db.playerDied(playerName)
		}
	}

	scoreChanged = db.createOrUpdateCategory(
		playerName,
		lineCatName,
		uint(playerCatRank),
		uint(playerCatScore),
	)

	return scoreChanged && playerIsAlive
}

/*
when dealing with api calls, only need to worry about new categories/ category updates, players won't be new
	also need to do a check for whether or not a player has died if we find an update before tweeting

when dealing with scraping, need to check if player exitsts

*/
