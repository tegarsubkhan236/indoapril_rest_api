package entities

import (
	"gorm.io/gorm"
	"time"
)

// CrTeam act as table, request body, response body
type CrTeam struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	ParentID  *uint          `json:"parent_id"`
	Name      string         `json:"name" gorm:"size:64;not null"`
	Children  []CrTeam       `json:"-" gorm:"foreignkey:ParentID"`
	Users     []CrUser       `json:"users" gorm:"foreignKey:TeamID"`
	CreatedAt time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
