package db

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DiscordID string `gorm:"uniqueIndex"`
	IGN       string
	Hash      string
}

type UserRepository struct {
	db *Db
}

func NewUserRepository(db *Db) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindAll() []User {
	var users []User
	r.db.Find(&users)

	return users
}

func (r *UserRepository) FindByDiscordID(discordID string) *User {
	var user User
	result := r.db.Where(&User{DiscordID: discordID}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &user
}

func (r *UserRepository) SetInfo(discordID, ign, hash string) *User {
	user := r.FindByDiscordID(discordID)
	if user != nil {
		user.IGN = ign
		user.Hash = hash
		r.db.Save(user)
	} else {
		user = &User{
			DiscordID: discordID,
			IGN:       ign,
			Hash:      hash,
		}
		r.db.Create(user)
	}

	return user
}
