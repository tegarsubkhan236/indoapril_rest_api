package model

import (
	"gorm.io/gorm"
	"time"
)

type TrPurchaseOrder struct {
	ID                      uint       `gorm:"primary_key"`
	MsSupplierID            uint       `gorm:"not null"`
	PoCode                  string     `gorm:"size:32;not null"`
	Disc                    int8       `gorm:"default:0"`
	Tax                     int8       `gorm:"default:0"`
	Amount                  int        `gorm:"default:0"`
	Remarks                 string     `gorm:"type:text"`
	Status                  int8       `gorm:"default:0"`
	MsSupplier              MsSupplier `gorm:"foreignkey:MsSupplierID"`
	TrPurchaseOrderProducts []TrPurchaseOrderProduct
	TrBackOrders            []TrBackOrder
	CreatedAt               time.Time      `gorm:"autoCreateTime"`
	UpdatedAt               time.Time      `gorm:"autoUpdateTime"`
	DeletedAt               gorm.DeletedAt `gorm:"index"`
}

type PurchaseOrderResponse struct {
	ID                    uint                           `json:"id"`
	PoCode                string                         `json:"po_code"`
	Disc                  int8                           `json:"disc"`
	Tax                   int8                           `json:"tax"`
	Amount                int                            `json:"amount"`
	Remarks               string                         `json:"remarks"`
	Status                int8                           `json:"status"`
	SupplierID            uint                           `json:"supplier_id"`
	SupplierName          string                         `json:"supplier_name"`
	PurchaseOrderProducts []PurchaseOrderProductResponse `json:"purchase_order_products"`
	BackOrders            []BackOrderResponse            `json:"back_orders"`
}

func ConvertPurchaseOrderToResponse(model TrPurchaseOrder) PurchaseOrderResponse {
	purchaseOrder := PurchaseOrderResponse{
		ID:           model.ID,
		PoCode:       model.PoCode,
		Disc:         model.Disc,
		Tax:          model.Tax,
		Amount:       model.Amount,
		Remarks:      model.Remarks,
		Status:       model.Status,
		SupplierID:   model.MsSupplier.ID,
		SupplierName: model.MsSupplier.Name,
	}
	for _, product := range model.TrPurchaseOrderProducts {
		item := ConvertPurchaseOrderProductToResponse(product)
		purchaseOrder.PurchaseOrderProducts = append(purchaseOrder.PurchaseOrderProducts, item)
	}
	for _, backOrder := range model.TrBackOrders {
		item := ConvertBackOrderToResponse(backOrder)
		purchaseOrder.BackOrders = append(purchaseOrder.BackOrders, item)
	}
	return purchaseOrder
}

func ConvertResponseToPurchaseOrder(response PurchaseOrderResponse) TrPurchaseOrder {
	purchaseOrder := TrPurchaseOrder{
		PoCode:       response.PoCode,
		MsSupplierID: response.SupplierID,
		Disc:         response.Disc,
		Tax:          response.Tax,
		Amount:       response.Amount,
		Remarks:      response.Remarks,
	}
	for _, product := range response.PurchaseOrderProducts {
		item := ConvertResponseToPurchaseOrderProduct(product)
		purchaseOrder.TrPurchaseOrderProducts = append(purchaseOrder.TrPurchaseOrderProducts, item)
	}
	return purchaseOrder
}
