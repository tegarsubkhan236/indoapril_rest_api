package entities

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

type TrPurchaseOrderReq struct {
	ID                    uint                        `json:"id"`
	PoCode                string                      `json:"po_code"`
	Disc                  int8                        `json:"disc"`
	Tax                   int8                        `json:"tax"`
	Amount                int                         `json:"amount"`
	Remarks               string                      `json:"remarks"`
	Status                int8                        `json:"status"`
	SupplierID            uint                        `json:"supplier_id"`
	SupplierName          string                      `json:"supplier_name"`
	PurchaseOrderProducts []TrPurchaseOrderProductReq `json:"purchase_order_products"`
	BackOrders            []TrBackOrderReq            `json:"back_orders"`
}

type TrPurchaseOrderResp struct {
	ID                    uint                         `json:"id"`
	PoCode                string                       `json:"po_code"`
	Disc                  int8                         `json:"disc"`
	Tax                   int8                         `json:"tax"`
	Amount                int                          `json:"amount"`
	Remarks               string                       `json:"remarks"`
	Status                int8                         `json:"status"`
	SupplierID            uint                         `json:"supplier_id"`
	SupplierName          string                       `json:"supplier_name"`
	PurchaseOrderProducts []TrPurchaseOrderProductResp `json:"purchase_order_products"`
	BackOrders            []TrBackOrderResp            `json:"back_orders"`
}

func (req TrPurchaseOrderReq) ToModel() TrPurchaseOrder {
	po := TrPurchaseOrder{
		PoCode:       req.PoCode,
		MsSupplierID: req.SupplierID,
		Disc:         req.Disc,
		Tax:          req.Tax,
		Amount:       req.Amount,
		Remarks:      req.Remarks,
	}
	for _, PoProduct := range req.PurchaseOrderProducts {
		item := TrPurchaseOrderProduct{
			MsProductID: PoProduct.ProductID,
			Price:       PoProduct.Price,
			Quantity:    PoProduct.Quantity,
		}
		po.TrPurchaseOrderProducts = append(po.TrPurchaseOrderProducts, item)
	}
	return po
}

func (model TrPurchaseOrder) ToResponse() TrPurchaseOrderResp {
	po := TrPurchaseOrderResp{
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
	for _, poProduct := range model.TrPurchaseOrderProducts {
		item := TrPurchaseOrderProductResp{
			ID:          poProduct.ID,
			ProductID:   poProduct.MsProduct.ID,
			ProductName: poProduct.MsProduct.Name,
			Quantity:    poProduct.Quantity,
			Price:       poProduct.Price,
		}
		po.PurchaseOrderProducts = append(po.PurchaseOrderProducts, item)
	}
	for _, backOrder := range model.TrBackOrders {
		po.BackOrders = append(po.BackOrders, backOrder.ToResponse())
	}

	return po
}
