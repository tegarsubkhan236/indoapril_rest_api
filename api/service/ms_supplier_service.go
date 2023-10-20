package service

import (
	"errors"
	"example/pkg"
	"example/pkg/model"
	"gorm.io/gorm"
)

func GetAllSupplier(page, limit int, ids []int, supplier model.MsSupplier) ([]model.MsSupplier, int64, error) {
	var db = pkg.DB
	var data []model.MsSupplier
	var count int64
	var offset = (page - 1) * limit

	db = db.Model(&model.MsSupplier{})

	if ids != nil {
		db = db.Where("id IN (?)", ids)
	}
	if supplier.Name != "" {
		db = db.Where("name LIKE ?", "%"+supplier.Name+"%")
	}
	if supplier.Address != "" {
		db = db.Where("address LIKE ?", "%"+supplier.Address+"%")
	}
	if supplier.ContactPerson != "" {
		db = db.Where("contact_person LIKE ?", "%"+supplier.ContactPerson+"%")
	}
	if supplier.ContactNumber != "" {
		db = db.Where("contact_number LIKE ?", "%"+supplier.ContactNumber+"%")
	}

	result := db.Count(&count).Offset(offset).Limit(limit).Order("id desc").Find(&data)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return data, count, nil
}

func GetSupplierById(id uint) (model.MsSupplier, error) {
	var db = pkg.DB
	var item model.MsSupplier

	if err := db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.MsSupplier{}, errors.New("item not found")
		}
		return model.MsSupplier{}, err
	}

	return item, nil
}

func CreateSupplier(data []model.MsSupplier) ([]model.MsSupplier, error) {
	var db = pkg.DB

	if err := db.Create(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func UpdateSupplier(item model.MsSupplier, payload model.MsSupplier) (model.MsSupplier, error) {
	var db = pkg.DB

	if err := db.Model(&item).Updates(payload).Error; err != nil {
		return model.MsSupplier{}, err
	}
	return item, nil
}

func DestroySupplier(ids []int) error {
	var db = pkg.DB
	var suppliers []model.MsSupplier
	var tx = db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := db.Unscoped().Delete(&suppliers, ids).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
