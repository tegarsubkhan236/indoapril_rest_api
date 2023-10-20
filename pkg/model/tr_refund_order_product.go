package model

import (
	"gorm.io/gorm"
	"time"
)

type TrReturnOrderProduct struct {
	ID              uint           `gorm:"primary_key"`
	Quantity        int            `gorm:"default:0"`
	Price           int            `gorm:"default:0"`
	TrReturnOrderID uint           `gorm:"not null"`
	TrReturnOrder   TrReturnOrder  `gorm:"foreignkey:TrReturnOrderID"`
	MsProductID     uint           `gorm:"not null"`
	MsProduct       MsProduct      `gorm:"foreignkey:MsProductID"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type ReturnOrderProductResponse struct {
	ID          uint   `json:"id,omitempty"`
	ProductID   uint   `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
	Price       int    `json:"price,omitempty"`
}

func ConvertReturnOrderProductToResponse(model TrReturnOrderProduct) ReturnOrderProductResponse {
	product := ReturnOrderProductResponse{
		ID:          model.ID,
		ProductID:   model.MsProduct.ID,
		ProductName: model.MsProduct.Name,
		Quantity:    model.Quantity,
		Price:       model.Price,
	}
	return product
}

func ConvertResponseToReturnOrderProduct(response ReturnOrderProductResponse) TrReturnOrderProduct {
	product := TrReturnOrderProduct{
		MsProductID: response.ProductID,
		Price:       response.Price,
		Quantity:    response.Quantity,
	}
	return product
}
