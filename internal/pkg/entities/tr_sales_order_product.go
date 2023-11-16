package entities

import (
	"gorm.io/gorm"
	"time"
)

type TrSalesOrderProduct struct {
	ID             uint           `gorm:"primary_key"`
	TrSalesOrderID uint           `gorm:"not null"`
	MsProductID    uint           `gorm:"not null"`
	Quantity       int            `gorm:"default:0"`
	TrSalesOrder   TrSalesOrder   `gorm:"foreignkey:TrSalesOrderID"`
	MsProduct      MsProduct      `gorm:"foreignkey:MsProductID"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type TrSalesOrderProductReq struct {
	ID        uint `json:"id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type TrSalesOrderProductResp struct {
	ID          uint   `json:"id"`
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}
