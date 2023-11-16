package ms_supplier

import (
	"errors"
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	CreateSupplier(data *[]entities.MsSupplier) (*[]entities.MsSupplier, error)
	ReadAllSupplier(page, limit int, ids []uint) (*[]entities.MsSupplier, int64, error)
	ReadSupplierById(id uint) (*entities.MsSupplier, error)
	UpdateSupplier(item *entities.MsSupplier, payload *entities.MsSupplier) (*entities.MsSupplier, error)
	DestroySupplier(ids []uint) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateSupplier(data *[]entities.MsSupplier) (*[]entities.MsSupplier, error) {
	var tx = r.DB.Begin()

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return data, nil
}

func (r repository) ReadAllSupplier(page, limit int, ids []uint) (*[]entities.MsSupplier, int64, error) {
	var data []entities.MsSupplier
	var count int64
	var offset = (page - 1) * limit

	r.DB = r.DB.Model(&entities.MsSupplier{})

	if len(ids) != 0 {
		r.DB = r.DB.Where("id IN (?)", ids)
	}

	if err := r.DB.Count(&count).Offset(offset).Limit(limit).Order("id desc").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return &data, count, nil
}

func (r repository) ReadSupplierById(id uint) (*entities.MsSupplier, error) {
	var item entities.MsSupplier

	if err := r.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}

	return &item, nil
}

func (r repository) UpdateSupplier(item *entities.MsSupplier, payload *entities.MsSupplier) (*entities.MsSupplier, error) {
	if err := r.DB.Model(&item).Updates(payload).Error; err != nil {
		return nil, err
	}

	return item, nil
}

func (r repository) DestroySupplier(ids []uint) error {
	var suppliers []entities.MsSupplier
	var tx = r.DB.Begin()

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Unscoped().Delete(&suppliers, ids).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
