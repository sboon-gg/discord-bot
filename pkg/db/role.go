package db

import (
	"errors"

	"gorm.io/gorm"
)

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

func (r *RoleRepository) FindByDiscordID(discordID string) *Role {
	var role Role
	result := r.db.Where(&Role{DiscordID: discordID}).First(&role)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return &role
}

func (r *RoleRepository) SetMapping(roleID, activeRoleID string) *Role {
	role := r.FindByDiscordID(roleID)
	if role != nil {
		role.ActiveRoleID = activeRoleID
		r.db.Save(role)
	} else {
		role = &Role{
			DiscordID:    roleID,
			ActiveRoleID: activeRoleID,
		}
		r.db.Create(role)
	}

	return role
}

func (r *RoleRepository) UnsetMapping(roleID string) {
	role := r.FindByDiscordID(roleID)
	if role != nil {
		r.db.Delete(role)
	}
}
