package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents an access role.
type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Menus     []RoleMenu     `gorm:"foreignKey:RoleID" json:"menus"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// RoleMenu stores which menu code is accessible by a role.
type RoleMenu struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	RoleID   uint   `gorm:"index;not null" json:"role_id"`
	MenuCode string `gorm:"size:100;not null" json:"menu_code"`
}
