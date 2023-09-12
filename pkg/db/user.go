package db

import "gorm.io/gorm"

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

func (r *UserRepository) FindByDiscordID(discordID string) User {
	var user User
	r.db.Where(&User{DiscordID: discordID}).First(&user)
	return user
}
