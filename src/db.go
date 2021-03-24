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

// ScrapeData struct used to determine whether or not to tweet, and what to tweet
type ChangeData struct {
	NewCategory   bool
	ScoreChanged  bool
	PreviousScore uint
	LastUpdate    time.Time
	PlayerAlive   bool
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

func (db *MyDB) createOrUpdateCategory(playerName string, catName string, playerRank uint, playerScore uint) *ChangeData {

	var category Category
	var changeData ChangeData

	newCategory := false
	scoreChanged := false

	playerID := db.getPlayerID(playerName)

	catDB := db.Table("categories").Where("player_id = ? AND name = ?", playerID, catName).First(&category)

	newCategory = errors.Is(catDB.Error, gorm.ErrRecordNotFound)

	row := db.Table("categories").Where("name = ? AND player_id = ?", catName, playerID).Select("score", "updated").Row()
	var score uint
	var updated time.Time
	row.Scan(&score, &updated)

	if newCategory {
		db.createCategory(
			playerName,
			catName,
			playerRank,
			playerScore,
		)
	} else {

		if score != playerScore {
			scoreChanged = true
		}
		db.updateCategory(playerName, catName, playerRank, playerScore)
	}

	changeData.ScoreChanged = scoreChanged
	changeData.NewCategory = newCategory
	changeData.PreviousScore = score
	changeData.LastUpdate = updated

	return &changeData
}

func (db *MyDB) highscoreLineCreateOrUpdate(highscoreLine *HighscoreLine) *ChangeData {

	newPlayer := false

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

	changeData := db.createOrUpdateCategory(
		playerName,
		lineCatName,
		uint(playerCatRank),
		uint(playerCatScore),
	)

	changeData.PlayerAlive = playerIsAlive

	return changeData
}

func (db *MyDB) getNextApiCallName() string {

	name := ""

	return name
}
