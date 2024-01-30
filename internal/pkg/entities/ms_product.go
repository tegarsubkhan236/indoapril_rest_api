package entities

import (
	"example/internal/pkg/types/stock_status"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type MsProduct struct {
	ID                  uint                `gorm:"primary_key"`
	MsSupplierID        uint                `gorm:"not null"`
	Name                string              `gorm:"size:32;not null"`
	Description         string              `gorm:"type:text"`
	Status              bool                `gorm:"default:0"`
	MsSupplier          MsSupplier          `gorm:"foreignkey:MsSupplierID"`
	MsProductCategories []MsProductCategory `gorm:"many2many:ms_product_product_categories"`
	MsProductPrices     []MsProductPrice    `json:"ms_product_prices"`
	MsStocks            []MsStock           `json:"ms_stocks"`
	CreatedAt           time.Time           `gorm:"autoCreateTime"`
	UpdatedAt           time.Time           `gorm:"autoUpdateTime"`
	DeletedAt           gorm.DeletedAt      `gorm:"index"`
}

type MsProductReq struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	SellPrice           int    `json:"sell_price"`
	UserID              uint   `json:"user_id"`
	SupplierID          uint   `json:"supplier_id"`
	ProductCategoriesID []uint `json:"product_categories_id"`
}

type MsProductResp struct {
	ID                uint                    `json:"id"`
	Name              string                  `json:"name"`
	Description       string                  `json:"description"`
	Status            bool                    `json:"status"`
	Supplier          MsSupplier              `json:"supplier"`
	ProductCategories []MsProductCategoryResp `json:"product_categories"`
	ProductPrices     []ProductPriceResp      `json:"product_prices"`
	Stocks            []StockResp             `json:"stocks"`
}

func (r MsProductReq) ValidateInput() error {
	if r.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if len(r.ProductCategoriesID) == 0 {
		return fmt.Errorf("product %s has no category", r.Name)
	}
	return nil
}

func (r MsProductReq) ToModel() MsProduct {
	item := MsProduct{
		Name:         r.Name,
		Description:  r.Description,
		Status:       true,
		MsSupplierID: r.SupplierID,
	}
	for _, categoryID := range r.ProductCategoriesID {
		x := MsProductCategory{
			ID: categoryID,
		}
		item.MsProductCategories = append(item.MsProductCategories, x)
	}
	item.MsProductPrices = []MsProductPrice{
		{
			SellPrice: r.SellPrice,
			BuyPrice:  0,
		},
	}
	item.MsStocks = []MsStock{
		{
			CoreUserID: r.UserID,
			Quantity:   0,
			Total:      0,
			Type:       stock_status.INITIAL_STOCK,
		},
	}
	return item
}

func (r MsProduct) ToResponse() MsProductResp {
	item := MsProductResp{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		Supplier: MsSupplier{
			ID:   r.MsSupplier.ID,
			Name: r.MsSupplier.Name,
		},
	}
	for _, productCategory := range r.MsProductCategories {
		x := MsProductCategoryResp{
			ID:       productCategory.ID,
			ParentID: productCategory.ParentID,
			Name:     productCategory.Name,
		}
		item.ProductCategories = append(item.ProductCategories, x)
	}
	for _, productPrice := range r.MsProductPrices {
		x := ProductPriceResp{
			ID:        productPrice.ID,
			BuyPrice:  productPrice.BuyPrice,
			SellPrice: productPrice.SellPrice,
		}
		item.ProductPrices = append(item.ProductPrices, x)
	}
	for _, productStock := range r.MsStocks {
		x := StockResp{
			ID:       productStock.ID,
			Quantity: productStock.Quantity,
			Total:    productStock.Total,
		}
		item.Stocks = append(item.Stocks, x)
	}
	return item
}
