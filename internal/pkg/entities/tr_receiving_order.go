package entities

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
	// TODO implement Created By UserID
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TrReceivingOrderReq struct {
	ID                     uint                         `json:"id"`
	UserID                 uint                         `json:"user_id"`
	PoID                   uint                         `json:"po_id"`
	BoID                   *uint                        `json:"bo_id"`
	PoCode                 string                       `json:"po_code"`
	BoCode                 string                       `json:"bo_code"`
	Amount                 int                          `json:"amount"`
	Remarks                string                       `json:"remarks"`
	ReceivingOrderProducts []TrReceivingOrderProductReq `json:"receiving_order_products"`
}

func (req TrReceivingOrderReq) ToModel() TrReceivingOrder {
	receivingOrder := TrReceivingOrder{
		TrPurchaseOrderID: req.PoID,
		TrBackOrderID:     req.BoID,
		Amount:            req.Amount,
		Remarks:           req.Remarks,
	}
	for _, val := range req.ReceivingOrderProducts {
		item := TrReceivingOrderProduct{
			MsProductID: val.ProductID,
			Quantity:    val.Quantity,
		}
		receivingOrder.TrReceivingProducts = append(receivingOrder.TrReceivingProducts, item)
	}
	return receivingOrder
}
