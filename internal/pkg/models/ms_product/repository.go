package ms_product

import (
	"errors"
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	CreateBatchProduct(payload *[]entities.MsProduct) (*[]entities.MsProduct, error)
	ReadAllProduct(page, limit int, productName string, supplierID uint, productCategories []uint) (*[]entities.MsProduct, int64, error)
	ReadProductById(id uint, basedOn string) (*entities.MsProduct, error)
	UpdateProduct(item *entities.MsProduct, payload *entities.MsProduct) (*entities.MsProduct, error)
	DestroyProduct(ids []uint) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateBatchProduct(payload *[]entities.MsProduct) (*[]entities.MsProduct, error) {
	if payload == nil && len(*payload) == 0 {
		return nil, errors.New("empty payload")
	}

	batchSize := 500

	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	for i := 0; i < len(*payload); i += batchSize {
		end := i + batchSize
		if end > len(*payload) {
			end = len(*payload)
		}
		batch := (*payload)[i:end]

		if err := tx.CreateInBatches(batch, len(batch)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return payload, nil
}

func (r repository) ReadAllProduct(page, limit int, productName string, supplierID uint, productCategories []uint) (*[]entities.MsProduct, int64, error) {
	var count int64
	var data []entities.MsProduct
	var offset = (page - 1) * limit

	r.DB = r.DB.Model(&data)

	if productName != "" {
		r.DB = r.DB.Where("ms_products.name LIKE ?", "%"+productName+"%")
	}

	// fetch supplier based on supplierID
	if supplierID != 0 {
		r.DB = r.DB.Joins("JOIN ms_suppliers ON ms_suppliers.id = ms_products.ms_supplier_id "+
			"AND ms_products.ms_supplier_id = ?", supplierID).Preload("MsSupplier")
	} else {
		r.DB = r.DB.Preload("MsSupplier")
	}

	// fetch product category based on productCategoryID
	if len(productCategories) != 0 {
		var ids []uint
		for _, item := range productCategories {
			ids = append(ids, item)
		}

		r.DB = r.DB.Joins("JOIN ms_product_product_categories ON ms_product_product_categories.ms_product_id = ms_products.id "+
			"JOIN ms_product_categories ON ms_product_product_categories.ms_product_category_id = ms_product_categories.id "+
			"AND ms_product_categories.id IN ? GROUP BY ms_products.id", ids).Preload("MsProductCategories")
	} else {
		r.DB = r.DB.Preload("MsProductCategories")
	}

	// fetch last product price
	r.DB = r.DB.Preload("MsProductPrices", func(db *gorm.DB) *gorm.DB {
		return db.Joins("INNER JOIN (SELECT ms_product_id, MAX(id) AS max_id FROM ms_product_prices GROUP BY ms_product_id) AS max_prices " +
			"ON ms_product_prices.ms_product_id = max_prices.ms_product_id AND ms_product_prices.id = max_prices.max_id")
	})

	// fetch last product stock
	r.DB = r.DB.Preload("MsStocks", func(db *gorm.DB) *gorm.DB {
		return db.Joins("INNER JOIN (SELECT ms_product_id, MAX(id) AS max_id FROM ms_stocks GROUP BY ms_product_id) AS max_stocks " +
			"ON ms_stocks.ms_product_id = max_stocks.ms_product_id AND ms_stocks.id = max_stocks.max_id")
	})

	if err := r.DB.Count(&count).Offset(offset).Limit(limit).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return &data, count, nil
}

func (r repository) ReadProductById(id uint, basedOn string) (*entities.MsProduct, error) {
	var item entities.MsProduct

	r.DB = r.DB.Model(&item)

	if basedOn == "price" {
		r.DB = r.DB.Preload("MsProductPrices")
	} else if basedOn == "stock" {
		r.DB = r.DB.Preload("MsStocks")
	} else {
		r.DB = r.DB.Preload("MsProductPrices", func(db *gorm.DB) *gorm.DB {
			return db.Joins("INNER JOIN (SELECT ms_product_id, MAX(id) AS max_id FROM ms_product_prices GROUP BY ms_product_id) AS max_prices " +
				"ON ms_product_prices.ms_product_id = max_prices.ms_product_id AND ms_product_prices.id = max_prices.max_id")
		})
		r.DB = r.DB.Preload("MsStocks", func(db *gorm.DB) *gorm.DB {
			return db.Joins("INNER JOIN (SELECT ms_product_id, MAX(id) AS max_id FROM ms_stocks GROUP BY ms_product_id) AS max_stocks " +
				"ON ms_stocks.ms_product_id = max_stocks.ms_product_id AND ms_stocks.id = max_stocks.max_id")
		})
	}

	if err := r.DB.Preload("MsSupplier").Preload("MsProductCategories").First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("item not found")
		}
		return nil, err
	}

	return &item, nil
}

func (r repository) UpdateProduct(item *entities.MsProduct, payload *entities.MsProduct) (*entities.MsProduct, error) {
	if err := r.DB.Model(&item).Updates(payload).Error; err != nil {
		return nil, err
	}
	return payload, nil
}

func (r repository) DestroyProduct(ids []uint) error {
	var products []entities.MsProduct
	var tx = r.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := r.DB.Where("id IN ?", ids).Find(&products).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, product := range products {
		if err := r.DB.Model(&product).Association("MsProductCategories").Clear(); err != nil {
			tx.Rollback()
			return err
		}

		if err := r.DB.Unscoped().Delete(&product).Error; err != nil {
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
