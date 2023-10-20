package service

import (
	"errors"
	"example/pkg"
	"example/pkg/model"
	"gorm.io/gorm"
)

func GetLastProductPrice(productID uint) (*model.MsProductPrice, error) {
	var lastMsProductPrice model.MsProductPrice
	var db = pkg.DB

	if err := db.Where("ms_product_id = ?", productID).Order("id desc").First(&lastMsProductPrice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("product price not found")
		}
		return nil, errors.New("failed to get product price: " + err.Error())
	}

	return &lastMsProductPrice, nil
}

func CreateProductPrice(ProductID uint, Price int) error {
	var db = pkg.DB
	var productPricePayload model.MsProductPrice

	lastMsProductPrice, err := GetLastProductPrice(ProductID)
	if err != nil {
		return err
	}
	if lastMsProductPrice.BuyPrice == Price {
		return nil
	}
	productPricePayload = model.MsProductPrice{
		MsProductID: lastMsProductPrice.MsProductID,
		SellPrice:   lastMsProductPrice.SellPrice,
		BuyPrice:    Price,
	}
	if err := db.Create(&productPricePayload).Error; err != nil {
		return err
	}
	return nil
}
