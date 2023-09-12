package db

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	DiscordID    string `gorm:"uniqueIndex"`
	ActiveRoleID string
}

type RoleRepository struct {
	db *Db
}

func NewRoleRepository(db *Db) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) FindAll() []Role {
	var roles []Role
	r.db.Find(&roles)
	return roles
}

func (r *RoleRepository) FindByDiscordID(discordID string) Role {
	var role Role
	r.db.Where(&Role{DiscordID: discordID}).First(&role)
	return role
}
