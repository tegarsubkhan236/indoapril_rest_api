package model

import (
	"gorm.io/gorm"
	"time"
)

type TrBackOrderProduct struct {
	ID            uint           `gorm:"primary_key"`
	TrBackOrderID uint           `gorm:"not null"`
	MsProductID   uint           `gorm:"not null"`
	Quantity      int            `gorm:"default:0"`
	Price         int            `gorm:"default:0"`
	TrBackOrder   TrBackOrder    `gorm:"foreignkey:TrBackOrderID"`
	MsProduct     MsProduct      `gorm:"foreignkey:MsProductID"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type BackOrderProductResponse struct {
	ID          uint   `json:"id,omitempty"`
	ProductID   uint   `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
	Price       int    `json:"price,omitempty"`
}
