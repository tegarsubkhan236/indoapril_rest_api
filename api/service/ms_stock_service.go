package service

import (
	"errors"
	"example/api/tool/converter"
	"example/pkg"
	"example/pkg/model"
	"gorm.io/gorm"
)

func GetLastStock(productID uint) (int, error) {
	var lastMsStock model.MsStock
	var db = pkg.DB

	err := db.Where("ms_product_id = ?", productID).Order("id desc").First(&lastMsStock).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	return lastMsStock.Total, nil
}

func CreateStock(typeStock int, userID, ProductID uint, Quantity int) error {
	var db = pkg.DB
	var stockPayload model.MsStock
	lastStock, errLastStock := GetLastStock(ProductID)
	if errLastStock != nil {
		return errLastStock
	}
	if typeStock == converter.Increment {
		stockPayload = model.MsStock{
			MsProductID: ProductID,
			CoreUserID:  userID,
			Quantity:    Quantity,
			Type:        converter.Increment,
			Total:       lastStock + Quantity,
		}
	}
	if typeStock == converter.Decrement {
		stockPayload = model.MsStock{
			MsProductID: ProductID,
			CoreUserID:  userID,
			Quantity:    Quantity,
			Type:        converter.Decrement,
			Total:       lastStock - Quantity,
		}
	}
	if errCreate := db.Create(&stockPayload).Error; errCreate != nil {
		return errCreate
	}
	return nil
}
