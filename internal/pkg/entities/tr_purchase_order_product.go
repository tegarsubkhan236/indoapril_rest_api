package entities

import (
	"gorm.io/gorm"
	"time"
)

type TrPurchaseOrderProduct struct {
	ID                uint            `gorm:"primary_key"`
	TrPurchaseOrderID uint            `gorm:"not null"`
	MsProductID       uint            `gorm:"not null"`
	Quantity          int             `gorm:"default:0"`
	Price             int             `gorm:"default:0"`
	TrPurchaseOrder   TrPurchaseOrder `gorm:"foreignkey:TrPurchaseOrderID"`
	MsProduct         MsProduct       `gorm:"foreignkey:MsProductID"`
	CreatedAt         time.Time       `gorm:"autoCreateTime"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt  `gorm:"index"`
}

type TrPurchaseOrderProductReq struct {
	ID        uint `json:"id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
	Price     int  `json:"price"`
}

type TrPurchaseOrderProductResp struct {
	ID          uint   `json:"id"`
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
}
