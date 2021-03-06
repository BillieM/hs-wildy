package main

import (
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

func (db *MyDB) createPlayer(name string, alive bool) {
	db.Create(&Player{
		Name:  name,
		Alive: alive,
	})
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

func (db *MyDB) createCategory(playerName string, catName string, rank uint, score uint) {
	playerID := db.getPlayerID(playerName)

	db.Create(&Category{
		Name:     catName,
		Rank:     rank,
		Score:    score,
		PlayerID: playerID,
	})
}

func (db *MyDB) updateCategory(playerName string, catName string, rank uint, score uint) {
	playerID := db.getPlayerID(playerName)

	db.Table("categories").Where("player_id = ? AND name = ?", playerID, catName).Updates(Category{
		Rank:  rank,
		Score: score,
	})

}

func (db *MyDB) createOrUpdateCategory(playerName string, catName string, rank uint, score uint) {
	/*
		tries to update category, failing that creates a new category
	*/
}
