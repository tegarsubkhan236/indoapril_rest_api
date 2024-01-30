package ms_stock

import (
	"errors"
	"example/internal/pkg/entities"
	"example/internal/pkg/types/stock_status"
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
	if typeStock != stock_status.INCREMENT && typeStock != stock_status.DECREMENT {
		return errors.New("invalid typeStock")
	}

	if Quantity <= 0 {
		return errors.New("quantity must be a positive integer")
	}

	var stockPayload entities.MsStock
	lastStock, err := r.ReadLastStock(ProductID)
	if err != nil {
		return err
	}

	if typeStock == stock_status.INCREMENT {
		stockPayload = entities.MsStock{
			MsProductID: ProductID,
			CoreUserID:  userID,
			Quantity:    Quantity,
			Type:        stock_status.INCREMENT,
			Total:       lastStock + Quantity,
		}
	}

	if typeStock == stock_status.DECREMENT {
		stockPayload = entities.MsStock{
			MsProductID: ProductID,
			CoreUserID:  userID,
			Quantity:    Quantity,
			Type:        stock_status.DECREMENT,
			Total:       lastStock - Quantity,
		}
	}

	if stockPayload.Total < 0 {
		return errors.New("the result of the total stock cannot be negative")
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
