package entities

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

type TrReceivingOrderProductReq struct {
	ID        uint `json:"id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
	Price     int  `json:"price"`
}

type TrReceivingOrderProductResp struct {
	ID        uint `json:"id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
	Price     int  `json:"price"`
}
