package ms_product_price

import (
	"errors"
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	CreateProductPrice(ProductID uint, Price int) error
	ReadLastProductPrice(productID uint) (*entities.MsProductPrice, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateProductPrice(ProductID uint, Price int) error {
	var productPricePayload entities.MsProductPrice

	lastMsProductPrice, err := r.ReadLastProductPrice(ProductID)
	if err != nil {
		return err
	}

	if lastMsProductPrice.BuyPrice == Price {
		return nil
	}

	productPricePayload = entities.MsProductPrice{
		MsProductID: lastMsProductPrice.MsProductID,
		SellPrice:   lastMsProductPrice.SellPrice,
		BuyPrice:    Price,
	}
	if err = r.DB.Create(&productPricePayload).Error; err != nil {
		return err
	}

	return nil
}

func (r repository) ReadLastProductPrice(productID uint) (*entities.MsProductPrice, error) {
	var lastMsProductPrice entities.MsProductPrice

	if err := r.DB.Where("ms_product_id = ?", productID).Order("id desc").First(&lastMsProductPrice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("product price not found")
		}
		return nil, errors.New("failed to get product price: " + err.Error())
	}

	return &lastMsProductPrice, nil
}