package model

import (
	"gorm.io/gorm"
	"time"
)

type TrReturnOrder struct {
	ID                    uint   `gorm:"primary_key"`
	RoCode                string `gorm:"size:32;not null"`
	Disc                  int8   `gorm:"default:0"`
	Tax                   int8   `gorm:"default:0"`
	Amount                int    `gorm:"default:0"`
	Remarks               string `gorm:"type:text"`
	Status                int8   `gorm:"default:0"`
	TrReturnOrderProducts []TrReturnOrderProduct
	CreatedAt             time.Time      `gorm:"autoCreateTime"`
	UpdatedAt             time.Time      `gorm:"autoUpdateTime"`
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}

type ReturnOrderResponse struct {
	ID                  uint                         `json:"id"`
	RoCode              string                       `json:"so_code"`
	Disc                int8                         `json:"disc"`
	Tax                 int8                         `json:"tax"`
	Amount              int                          `json:"amount"`
	Remarks             string                       `json:"remarks"`
	Status              int8                         `json:"status"`
	ReturnOrderProducts []ReturnOrderProductResponse `json:"return_order_products"`
}

func ConvertReturnOrderToResponse(model TrReturnOrder) ReturnOrderResponse {
	salesOrder := ReturnOrderResponse{
		ID:      model.ID,
		RoCode:  model.RoCode,
		Disc:    model.Disc,
		Tax:     model.Tax,
		Amount:  model.Amount,
		Remarks: model.Remarks,
		Status:  model.Status,
	}
	for _, product := range model.TrReturnOrderProducts {
		item := ConvertReturnOrderProductToResponse(product)
		salesOrder.ReturnOrderProducts = append(salesOrder.ReturnOrderProducts, item)
	}
	return salesOrder
}

func ConvertResponseToReturnOrder(response ReturnOrderResponse) TrReturnOrder {
	salesOrder := TrReturnOrder{
		RoCode:  response.RoCode,
		Disc:    response.Disc,
		Tax:     response.Tax,
		Amount:  response.Amount,
		Remarks: response.Remarks,
	}
	for _, product := range response.ReturnOrderProducts {
		item := ConvertResponseToReturnOrderProduct(product)
		salesOrder.TrReturnOrderProducts = append(salesOrder.TrReturnOrderProducts, item)
	}
	return salesOrder
}
