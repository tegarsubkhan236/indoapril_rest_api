package model

import (
	"gorm.io/gorm"
	"time"
)

type TrBackOrder struct {
	ID                  uint            `gorm:"primary_key"`
	TrPurchaseOrderID   uint            `gorm:"not null"`
	MsSupplierID        uint            `gorm:"not null"`
	BoCode              string          `gorm:"size:32;not null"`
	Disc                int8            `gorm:"default:0"`
	Tax                 int8            `gorm:"default:0"`
	Amount              int             `gorm:"default:0"`
	Remarks             string          `gorm:"type:text"`
	Status              int8            `gorm:"default:0"`
	TrPurchaseOrder     TrPurchaseOrder `gorm:"foreignkey:TrPurchaseOrderID"`
	MsSupplier          MsSupplier      `gorm:"foreignkey:MsSupplierID"`
	TrBackOrderProducts []TrBackOrderProduct
	CreatedAt           time.Time      `gorm:"autoCreateTime"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

type BackOrderResponse struct {
	ID                uint                       `json:"id"`
	BoCode            string                     `json:"bo_code"`
	Disc              int8                       `json:"disc"`
	Tax               int8                       `json:"tax"`
	Amount            int                        `json:"amount"`
	Remarks           string                     `json:"remarks"`
	Status            int8                       `json:"status"`
	SupplierID        uint                       `json:"supplier_id"`
	SupplierName      string                     `json:"supplier_name"`
	BackOrderProducts []BackOrderProductResponse `json:"back_order_products"`
}

func ConvertBackOrderToResponse(backOrder TrBackOrder) BackOrderResponse {
	item := BackOrderResponse{
		ID:           backOrder.ID,
		BoCode:       backOrder.BoCode,
		Disc:         backOrder.Disc,
		Tax:          backOrder.Tax,
		Amount:       backOrder.Amount,
		Remarks:      backOrder.Remarks,
		Status:       backOrder.Status,
		SupplierID:   backOrder.MsSupplier.ID,
		SupplierName: backOrder.MsSupplier.Name,
	}
	for _, product := range backOrder.TrBackOrderProducts {
		backOrderProduct := BackOrderProductResponse{
			ID:          product.ID,
			ProductID:   product.MsProduct.ID,
			ProductName: product.MsProduct.Name,
			Quantity:    product.Quantity,
			Price:       product.Price,
		}
		item.BackOrderProducts = append(item.BackOrderProducts, backOrderProduct)
	}
	return item
}
