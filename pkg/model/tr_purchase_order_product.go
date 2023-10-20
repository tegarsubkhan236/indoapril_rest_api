package model

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

type PurchaseOrderProductResponse struct {
	ID          uint   `json:"id,omitempty"`
	ProductID   uint   `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
	Price       int    `json:"price,omitempty"`
}

func ConvertPurchaseOrderProductToResponse(model TrPurchaseOrderProduct) PurchaseOrderProductResponse {
	product := PurchaseOrderProductResponse{
		ID:          model.ID,
		ProductID:   model.MsProduct.ID,
		ProductName: model.MsProduct.Name,
		Quantity:    model.Quantity,
		Price:       model.Price,
	}
	return product
}

func ConvertResponseToPurchaseOrderProduct(response PurchaseOrderProductResponse) TrPurchaseOrderProduct {
	product := TrPurchaseOrderProduct{
		MsProductID: response.ProductID,
		Price:       response.Price,
		Quantity:    response.Quantity,
	}
	return product
}
