package entities

import (
	"gorm.io/gorm"
	"time"
)

// CrTeam act as table and request body
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

// CrTeamResp act as response body
type CrTeamResp struct {
	ID       uint         `json:"id"`
	ParentID *uint        `json:"parent_id,omitempty"`
	Name     string       `json:"name"`
	Children []CrTeamResp `json:"children,omitempty"`
	Users    []CrUserResp `json:"users,omitempty"`
}
