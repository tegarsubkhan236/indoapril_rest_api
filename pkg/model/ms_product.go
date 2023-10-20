package model

import (
	"example/api/tool/converter"
	"gorm.io/gorm"
	"time"
)

type MsProduct struct {
	ID                  uint                `gorm:"primary_key"`
	MsSupplierID        uint                `gorm:"not null"`
	Name                string              `gorm:"size:32;not null"`
	Description         string              `gorm:"type:text"`
	Status              int8                `gorm:"default:0"`
	MsSupplier          MsSupplier          `gorm:"foreignkey:MsSupplierID"`
	MsProductCategories []MsProductCategory `gorm:"many2many:ms_product_product_categories"`
	MsProductPrices     []MsProductPrice
	MsStocks            []MsStock
	CreatedAt           time.Time      `gorm:"autoCreateTime"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

type ProductRequest struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	SellPrice           int    `json:"sell_price"`
	UserID              uint   `json:"user_id"`
	SupplierID          uint   `json:"supplier_id"`
	ProductCategoriesID []uint `json:"product_categories_id"`
}

type ProductResponse struct {
	ID                uint                   `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Status            int8                   `json:"status"`
	Supplier          MsSupplier             `json:"supplier"`
	ProductCategories []MsProductCategory    `json:"product_categories"`
	ProductPrices     []ProductPriceResponse `json:"product_prices"`
	Stocks            []StockResponse        `json:"stocks"`
}

func (model MsProduct) ToResponse() ProductResponse {
	item := ProductResponse{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Status:      model.Status,
		Supplier: MsSupplier{
			ID:   model.MsSupplier.ID,
			Name: model.MsSupplier.Name,
		},
	}
	for _, productCategory := range model.MsProductCategories {
		x := MsProductCategory{
			ID:       productCategory.ID,
			ParentID: productCategory.ParentID,
			Name:     productCategory.Name,
		}
		item.ProductCategories = append(item.ProductCategories, x)
	}
	for _, productPrice := range model.MsProductPrices {
		x := ProductPriceResponse{
			ID:        productPrice.ID,
			BuyPrice:  productPrice.BuyPrice,
			SellPrice: productPrice.SellPrice,
		}
		item.ProductPrices = append(item.ProductPrices, x)
	}
	for _, productStock := range model.MsStocks {
		x := StockResponse{
			ID:       productStock.ID,
			Quantity: productStock.Quantity,
			Total:    productStock.Total,
		}
		item.Stocks = append(item.Stocks, x)
	}
	return item
}

func (req ProductRequest) ToModel() MsProduct {
	product := MsProduct{
		Name:         req.Name,
		Description:  req.Description,
		MsSupplierID: req.SupplierID,
	}
	for _, categoryID := range req.ProductCategoriesID {
		category := MsProductCategory{ID: categoryID}
		product.MsProductCategories = append(product.MsProductCategories, category)
	}
	product.MsProductPrices = []MsProductPrice{
		{
			SellPrice: req.SellPrice,
			BuyPrice:  0,
		},
	}
	product.MsStocks = []MsStock{
		{
			CoreUserID: req.UserID,
			Quantity:   converter.InitialStock,
			Total:      converter.InitialStock,
			Type:       converter.InitialStock,
		},
	}
	return product
}
