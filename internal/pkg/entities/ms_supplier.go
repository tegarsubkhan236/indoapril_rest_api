package entities

import (
	"gorm.io/gorm"
	"time"
)

type MsSupplier struct {
	ID            uint           `gorm:"primary_key" json:"id"`
	Name          string         `gorm:"size:32;not null" json:"name"`
	Address       string         `gorm:"size:150;not null" json:"address"`
	ContactPerson string         `gorm:"size:32;not null" json:"contact_person"`
	ContactNumber string         `gorm:"size:32;not null" json:"contact_number"`
	Status        int8           `gorm:"default:0" json:"status"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
