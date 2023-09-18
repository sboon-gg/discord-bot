package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func New() *Db {
	db, err := gorm.Open(sqlite.Open("bot.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Role{})

	return &Db{
		db,
	}
}
