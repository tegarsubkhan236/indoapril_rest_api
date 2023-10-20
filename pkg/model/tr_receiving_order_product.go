package model

import (
	"gorm.io/gorm"
	"time"
)

type TrReceivingOrderProduct struct {
	ID                 uint             `gorm:"primary_key"`
	TrReceivingOrderID uint             `gorm:"not null"`
	MsProductID        uint             `gorm:"not null"`
	Quantity           int              `gorm:"default:0"`
	TrReceivingOrder   TrReceivingOrder `gorm:"foreignkey:TrReceivingOrderID"`
	MsProduct          MsProduct        `gorm:"foreignkey:MsProductID"`
	CreatedAt          time.Time        `gorm:"autoCreateTime"`
	UpdatedAt          time.Time        `gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt   `gorm:"index"`
}

type ReceivingOrderProductResponse struct {
	ID        uint `json:"id,omitempty"`
	ProductID uint `json:"product_id,omitempty"`
	Quantity  int  `json:"quantity,omitempty"`
	Price     int  `json:"price,omitempty"`
}
