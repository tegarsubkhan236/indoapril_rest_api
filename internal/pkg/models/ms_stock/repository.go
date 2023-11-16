package ms_stock

import (
	"errors"
	"example/internal/pkg/entities"
	"example/internal/pkg/util"
	"gorm.io/gorm"
)

type Repository interface {
	CreateStock(typeStock int, userID, ProductID uint, Quantity int) error
	ReadLastStock(productID uint) (int, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateStock(typeStock int, userID, ProductID uint, Quantity int) error {
	var stockPayload entities.MsStock
	lastStock, errLastStock := r.ReadLastStock(ProductID)
	if errLastStock != nil {
		return errLastStock
	}
	if typeStock == util.Increment {
		stockPayload = entities.MsStock{
			MsProductID: ProductID,
			CoreUserID:  userID,
			Quantity:    Quantity,
			Type:        util.Increment,
			Total:       lastStock + Quantity,
		}
	}
	if typeStock == util.Decrement {
		stockPayload = entities.MsStock{
			MsProductID: ProductID,
			CoreUserID:  userID,
			Quantity:    Quantity,
			Type:        util.Decrement,
			Total:       lastStock - Quantity,
		}
	}
	if errCreate := r.DB.Create(&stockPayload).Error; errCreate != nil {
		return errCreate
	}
	return nil
}

func (r repository) ReadLastStock(productID uint) (int, error) {
	var lastMsStock entities.MsStock

	err := r.DB.Where("ms_product_id = ?", productID).Order("id desc").First(&lastMsStock).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	return lastMsStock.Total, nil
}