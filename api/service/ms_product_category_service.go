package service

import (
	"errors"
	"example/pkg"
	"example/pkg/model"
	"gorm.io/gorm"
)

func GetAllProductCategory(offset, limit int) ([]model.MsProductCategory, int64, error) {
	var db = pkg.DB
	var count int64
	var data []model.MsProductCategory

	db = db.Model(&data)
	db = db.Preload("Children").Preload("Children.Children")

	if err := db.Count(&count).Where("parent_id IS NULL").Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, errors.New("failed to get product categories: " + err.Error())
	}

	return data, count, nil
}

func GetProductCategoryById(id string) (*model.MsProductCategory, error) {
	var db = pkg.DB
	var item model.MsProductCategory

	db = db.Preload("Children").Preload("Children.Children")

	if err := db.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("product category not found")
		}
		return nil, errors.New("failed to get product category: " + err.Error())
	}

	return &item, nil
}

func CreateProductCategory(data model.MsProductCategory) error {
	var db = pkg.DB

	if err := db.Create(&data).Error; err != nil {
		return errors.New("failed to create product category: " + err.Error())
	}

	return nil
}

func UpdateProductCategory(item model.MsProductCategory, payload model.MsProductCategory) error {
	var db = pkg.DB

	if err := db.Model(&item).Updates(payload).Error; err != nil {
		return errors.New("failed to update product category: " + err.Error())
	}

	return nil
}

func DestroyProductCategory(id string) error {
	var db = pkg.DB
	var item model.MsProductCategory

	if err := db.Unscoped().Delete(&item, "id = ?", id).Error; err != nil {
		return errors.New("failed to delete product category: " + err.Error())
	}

	return nil
}
