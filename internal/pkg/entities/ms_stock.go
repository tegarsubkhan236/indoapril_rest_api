package entities

import (
	"gorm.io/gorm"
	"time"
)

type MsStock struct {
	ID          uint `gorm:"primary_key"`
	MsProductID uint
	CoreUserID  uint
	Quantity    int            `gorm:"default:0"`
	Total       int            `gorm:"default:0"`
	Type        int8           `gorm:"not null"`
	MsProduct   MsProduct      `gorm:"foreignkey:MsProductID"`
	CoreUser    CrUser         `gorm:"foreignkey:CoreUserID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type StockResp struct {
	ID       uint   `json:"id"`
	UserID   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Quantity int    `json:"quantity"`
	Total    int    `json:"total"`
	Type     int8   `json:"type"`
}
