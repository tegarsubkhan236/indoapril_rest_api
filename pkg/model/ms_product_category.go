package model

import (
	"gorm.io/gorm"
	"time"
)

type MsProductCategory struct {
	ID        uint                `gorm:"primary_key"`
	ParentID  *uint               `json:"parent_id"`
	Name      string              `gorm:"size:32;not null" json:"name"`
	Children  []MsProductCategory `gorm:"foreignkey:ParentID" json:"children,omitempty"`
	CreatedAt time.Time           `gorm:"autoCreateTime"`
	UpdatedAt time.Time           `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt      `gorm:"index"`
}
