package model

import (
	"gorm.io/gorm"
	"time"
)

type TrReceivingOrder struct {
	ID                  uint `gorm:"primary_key"`
	TrPurchaseOrderID   uint
	TrBackOrderID       *uint
	Amount              int    `gorm:"default:0"`
	Remarks             string `gorm:"type:text"`
	TrPurchaseOrder     TrPurchaseOrder
	TrBackOrder         TrBackOrder
	TrReceivingProducts []TrReceivingOrderProduct
	CreatedAt           time.Time      `gorm:"autoCreateTime"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

type ReceivingOrderResponse struct {
	ID                     uint                            `json:"id"`
	UserID                 uint                            `json:"user_id"`
	PoID                   uint                            `json:"po_id"`
	PoCode                 string                          `json:"po_code"`
	BoID                   uint                            `json:"bo_id"`
	BoCode                 string                          `json:"bo_code"`
	Amount                 int                             `json:"amount"`
	Remarks                string                          `json:"remarks"`
	ReceivingOrderProducts []ReceivingOrderProductResponse `json:"receiving_order_products"`
}

func ConvertResponseToReceivingOrder(resp ReceivingOrderResponse) TrReceivingOrder {
	var boID *uint
	if resp.BoID != 0 {
		boID = &resp.BoID
	}
	receivingOrder := TrReceivingOrder{
		TrPurchaseOrderID: resp.PoID,
		TrBackOrderID:     boID,
		Amount:            resp.Amount,
		Remarks:           resp.Remarks,
	}
	for _, val := range resp.ReceivingOrderProducts {
		item := TrReceivingOrderProduct{
			MsProductID: val.ProductID,
			Quantity:    val.Quantity,
		}
		receivingOrder.TrReceivingProducts = append(receivingOrder.TrReceivingProducts, item)
	}
	return receivingOrder
}
