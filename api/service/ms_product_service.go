package service

import (
	"errors"
	"example/pkg"
	"example/pkg/model"
	"gorm.io/gorm"
)

func GetAllProduct(page, limit, supplierID int, productCategoryID []int, searchText string) ([]model.MsProduct, int64, error) {
	var db = pkg.DB
	var count int64
	var data []model.MsProduct
	var offset = (page - 1) * limit

	db = db.Model(&data)

	// fetch supplier based on supplierID
	if supplierID != 0 {
		db = db.Joins("JOIN ms_suppliers ON ms_suppliers.id = ms_products.ms_supplier_id "+
			"AND ms_products.ms_supplier_id = ?", supplierID).Preload("MsSupplier")
	} else {
		db = db.Preload("MsSupplier")
	}

	// fetch product category based on productCategoryID
	if len(productCategoryID) != 0 {
		db = db.Joins("JOIN ms_product_product_categories ON ms_product_product_categories.ms_product_id = ms_products.id "+
			"JOIN ms_product_categories ON ms_product_product_categories.ms_product_category_id = ms_product_categories.id "+
			"AND ms_product_categories.id IN ? GROUP BY ms_products.id", productCategoryID).Preload("MsProductCategories")
	} else {
		db = db.Preload("MsProductCategories")
	}

	// fetch last product price
	db = db.Preload("MsProductPrices", func(db *gorm.DB) *gorm.DB {
		return db.Joins("INNER JOIN (SELECT ms_product_id, MAX(id) AS max_id FROM ms_product_prices GROUP BY ms_product_id) AS max_prices " +
			"ON ms_product_prices.ms_product_id = max_prices.ms_product_id AND ms_product_prices.id = max_prices.max_id")
	})

	// fetch last product stock
	db = db.Preload("MsStocks", func(db *gorm.DB) *gorm.DB {
		return db.Joins("INNER JOIN (SELECT ms_product_id, MAX(id) AS max_id FROM ms_stocks GROUP BY ms_product_id) AS max_stocks " +
			"ON ms_stocks.ms_product_id = max_stocks.ms_product_id AND ms_stocks.id = max_stocks.max_id")
	})

	if searchText != "" {
		db = db.Where("ms_products.name LIKE ?", "%"+searchText+"%")
	}

	if err := db.Count(&count).Offset(offset).Limit(limit).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func GetProductById(id uint, basedOn string) (*model.MsProduct, error) {
	var db = pkg.DB
	var item model.MsProduct

	db = db.Model(&item)

	if basedOn == "price" {
		db = db.Preload("MsProductPrices")
	} else if basedOn == "stock" {
		db = db.Preload("MsStocks")
	} else {
		db = db.Preload("MsProductPrices", func(db *gorm.DB) *gorm.DB {
			return db.Joins("INNER JOIN (SELECT ms_product_id, MAX(id) AS max_id FROM ms_product_prices GROUP BY ms_product_id) AS max_prices " +
				"ON ms_product_prices.ms_product_id = max_prices.ms_product_id AND ms_product_prices.id = max_prices.max_id")
		})
		db = db.Preload("MsStocks", func(db *gorm.DB) *gorm.DB {
			return db.Joins("INNER JOIN (SELECT ms_product_id, MAX(id) AS max_id FROM ms_stocks GROUP BY ms_product_id) AS max_stocks " +
				"ON ms_stocks.ms_product_id = max_stocks.ms_product_id AND ms_stocks.id = max_stocks.max_id")
		})
	}

	if err := db.Preload("MsSupplier").Preload("MsProductCategories").First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("item not found")
		}
		return nil, err
	}

	return &item, nil
}

func CreateProduct(data []model.MsProduct) ([]model.MsProduct, error) {
	var db = pkg.DB
	batchSize := 500

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := data[i:end]

		if err := tx.CreateInBatches(batch, len(batch)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return data, nil
}

func UpdateProduct(item model.MsProduct, payload model.MsProduct) (*model.MsProduct, error) {
	var db = pkg.DB

	if err := db.Model(&item).Updates(payload).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func DestroyProduct(ids []int) error {
	var db = pkg.DB
	var products []model.MsProduct
	var tx = db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := db.Where("id IN ?", ids).Find(&products).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, product := range products {
		if err := db.Model(&product).Association("MsProductCategories").Clear(); err != nil {
			tx.Rollback()
			return err
		}

		if err := db.Unscoped().Delete(&product).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
