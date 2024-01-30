package entities

import (
	"gorm.io/gorm"
	"time"
)

type MsProductPrice struct {
	ID          uint           `json:"id" gorm:"primary_key"`
	MsProductID uint           `json:"-" gorm:"not null"`
	SellPrice   int            `json:"sell_price" gorm:"default:0"`
	BuyPrice    int            `json:"buy_price" gorm:"default:0"`
	MsProduct   MsProduct      `json:"-" gorm:"foreignkey:MsProductID"`
	CreatedAt   time.Time      `json:"-" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type ProductPriceResp struct {
	ID        uint `json:"id"`
	SellPrice int  `json:"sell_price"`
	BuyPrice  int  `json:"buy_price"`
}
