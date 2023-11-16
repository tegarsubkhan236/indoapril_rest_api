package ms_product_category

import (
	"errors"
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	CreateProductCategory(data *entities.MsProductCategory) error
	ReadAllProductCategory(offset, limit int) (*[]entities.MsProductCategory, int64, error)
	ReadProductCategoryById(id uint) (*entities.MsProductCategory, error)
	UpdateProductCategory(item *entities.MsProductCategory, payload *entities.MsProductCategory) error
	DestroyProductCategory(id uint) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateProductCategory(data *entities.MsProductCategory) error {
	if err := r.DB.Create(&data).Error; err != nil {
		return errors.New("failed to create product category: " + err.Error())
	}

	return nil
}

func (r repository) ReadAllProductCategory(offset, limit int) (*[]entities.MsProductCategory, int64, error) {
	var count int64
	var data []entities.MsProductCategory

	r.DB = r.DB.Model(&data)
	r.DB = r.DB.Preload("Children").Preload("Children.Children")

	if err := r.DB.Count(&count).Where("parent_id IS NULL").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, errors.New("failed to get product categories: " + err.Error())
	}

	return &data, count, nil
}

func (r repository) ReadProductCategoryById(id uint) (*entities.MsProductCategory, error) {
	var item entities.MsProductCategory

	r.DB = r.DB.Preload("Children").Preload("Children.Children")

	if err := r.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("product category not found")
		}
		return nil, errors.New("failed to get product category: " + err.Error())
	}

	return &item, nil
}

func (r repository) UpdateProductCategory(item *entities.MsProductCategory, payload *entities.MsProductCategory) error {
	if err := r.DB.Model(&item).Updates(payload).Error; err != nil {
		return errors.New("failed to update product category: " + err.Error())
	}

	return nil
}

func (r repository) DestroyProductCategory(id uint) error {
	var item entities.MsProductCategory

	if err := r.DB.Unscoped().Delete(&item, "id = ?", id).Error; err != nil {
		return errors.New("failed to delete product category: " + err.Error())
	}

	return nil
}
