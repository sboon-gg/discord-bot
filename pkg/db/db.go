package db

import (
	"github.com/glebarez/sqlite"
	"github.com/sboon-gg/sboon-bot/pkg/config"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func New(conf *config.DbConfig) *Db {
	db, err := gorm.Open(sqlite.Open(conf.Filename), &gorm.Config{})
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
