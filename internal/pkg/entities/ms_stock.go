package entities

import (
	"gorm.io/gorm"
	"time"
)

type MsStock struct {
	ID          uint           `json:"id" gorm:"primary_key"`
	MsProductID uint           `json:"-"`
	CoreUserID  uint           `json:"-"`
	Quantity    int            `json:"quantity" gorm:"default:0"`
	Total       int            `json:"total" gorm:"default:0"`
	Type        int8           `json:"type" gorm:"not null"`
	MsProduct   MsProduct      `json:"-" gorm:"foreignkey:MsProductID"`
	CoreUser    CrUser         `json:"-" gorm:"foreignkey:CoreUserID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt   time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type StockResp struct {
	ID       uint   `json:"id"`
	UserID   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Quantity int    `json:"quantity"`
	Total    int    `json:"total"`
	Type     int8   `json:"type"`
}
