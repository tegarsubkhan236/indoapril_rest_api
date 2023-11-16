package entities

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

type TrSalesOrderReq struct {
	SoCode             string                   `json:"so_code"`
	UserID             uint                     `json:"user_id"`
	Disc               int8                     `json:"disc"`
	Tax                int8                     `json:"tax"`
	Amount             int                      `json:"amount"`
	Remarks            string                   `json:"remarks"`
	SalesOrderProducts []TrSalesOrderProductReq `json:"sales_order_products"`
}

type TrSalesOrderResp struct {
	ID                 uint                      `json:"id"`
	SoCode             string                    `json:"so_code"`
	UserID             uint                      `json:"user_id"`
	Disc               int8                      `json:"disc"`
	Tax                int8                      `json:"tax"`
	Amount             int                       `json:"amount"`
	Remarks            string                    `json:"remarks"`
	Status             int8                      `json:"status"`
	SalesOrderProducts []TrSalesOrderProductResp `json:"sales_order_products"`
}

func (req TrSalesOrderReq) ToModel() TrSalesOrder {
	item := TrSalesOrder{
		SoCode:  req.SoCode,
		Disc:    req.Disc,
		Tax:     req.Tax,
		Amount:  req.Amount,
		Remarks: req.Remarks,
	}

	for _, x := range req.SalesOrderProducts {
		product := TrSalesOrderProduct{
			MsProductID: x.ProductID,
			Quantity:    x.Quantity,
		}
		item.TrSalesOrderProducts = append(item.TrSalesOrderProducts, product)
	}

	return item
}

func (s TrSalesOrder) ToResponse() TrSalesOrderResp {
	item := TrSalesOrderResp{
		SoCode:  s.SoCode,
		Disc:    s.Disc,
		Tax:     s.Tax,
		Amount:  s.Amount,
		Remarks: s.Remarks,
	}

	for _, x := range s.TrSalesOrderProducts {
		product := TrSalesOrderProductResp{
			ProductID:   x.MsProduct.ID,
			ProductName: x.MsProduct.Name,
			Quantity:    x.Quantity,
		}

		item.SalesOrderProducts = append(item.SalesOrderProducts, product)
	}

	return item
}
