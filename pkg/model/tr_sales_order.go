package model

import (
	"gorm.io/gorm"
	"time"
)

type TrSalesOrder struct {
	ID                   uint   `gorm:"primary_key"`
	SoCode               string `gorm:"size:32;not null"`
	Disc                 int8   `gorm:"default:0"`
	Tax                  int8   `gorm:"default:0"`
	Amount               int    `gorm:"default:0"`
	Remarks              string `gorm:"type:text"`
	Status               int8   `gorm:"default:0"`
	TrSalesOrderProducts []TrSalesOrderProduct
	CreatedAt            time.Time      `gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

type SalesOrderResponse struct {
	ID                 uint                        `json:"id"`
	SoCode             string                      `json:"so_code"`
	UserID             uint                        `json:"user_id"`
	Disc               int8                        `json:"disc"`
	Tax                int8                        `json:"tax"`
	Amount             int                         `json:"amount"`
	Remarks            string                      `json:"remarks"`
	Status             int8                        `json:"status"`
	SalesOrderProducts []SalesOrderProductResponse `json:"sales_order_products"`
}

func ConvertSalesOrderToResponse(model TrSalesOrder) SalesOrderResponse {
	salesOrder := SalesOrderResponse{
		ID:      model.ID,
		SoCode:  model.SoCode,
		Disc:    model.Disc,
		Tax:     model.Tax,
		Amount:  model.Amount,
		Remarks: model.Remarks,
		Status:  model.Status,
	}
	for _, product := range model.TrSalesOrderProducts {
		item := ConvertSalesOrderProductToResponse(product)
		salesOrder.SalesOrderProducts = append(salesOrder.SalesOrderProducts, item)
	}
	return salesOrder
}

func ConvertResponseToSalesOrder(response SalesOrderResponse) TrSalesOrder {
	salesOrder := TrSalesOrder{
		SoCode:  response.SoCode,
		Disc:    response.Disc,
		Tax:     response.Tax,
		Amount:  response.Amount,
		Remarks: response.Remarks,
	}
	for _, product := range response.SalesOrderProducts {
		item := ConvertResponseToSalesOrderProduct(product)
		salesOrder.TrSalesOrderProducts = append(salesOrder.TrSalesOrderProducts, item)
	}
	return salesOrder
}
