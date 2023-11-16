package entities

import (
	"gorm.io/gorm"
	"time"
)

type MsProductPrice struct {
	ID          uint           `gorm:"primary_key"`
	MsProductID uint           `gorm:"not null"`
	SellPrice   int            `gorm:"default:0"`
	BuyPrice    int            `gorm:"default:0"`
	MsProduct   MsProduct      `gorm:"foreignkey:MsProductID"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type ProductPriceResp struct {
	ID        uint `json:"id"`
	SellPrice int  `json:"sell_price"`
	BuyPrice  int  `json:"buy_price"`
}
