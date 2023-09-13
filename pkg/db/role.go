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

func (r *RoleRepository) SetMapping(roleID, activeRoleID string) (*Role, error) {
	role := r.FindByDiscordID(roleID)
	if role != nil {
		role.ActiveRoleID = activeRoleID

		result := r.db.Save(role)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		role = &Role{
			DiscordID:    roleID,
			ActiveRoleID: activeRoleID,
		}

		result := r.db.Create(role)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return role, nil
}

func (r *RoleRepository) UnsetMapping(roleID string) error {
	role := r.FindByDiscordID(roleID)
	if role != nil {
		result := r.db.Delete(role)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
