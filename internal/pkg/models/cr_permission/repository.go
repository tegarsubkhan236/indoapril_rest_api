package cr_permission

import (
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	Create(payload *entities.CrPermission) error
	ReadAll(page, limit int) (*[]entities.CrPermission, int64, error)
	ReadByID(ID uint) (*entities.CrPermission, error)
	Update(item *entities.CrPermission, payload *entities.CrPermission) error
	Delete(ID uint) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) Create(payload *entities.CrPermission) error {
	if err := r.DB.Create(&payload).Error; err != nil {
		return err
	}

	return nil
}

func (r repository) ReadAll(page, limit int) (*[]entities.CrPermission, int64, error) {
	var count int64
	var data []entities.CrPermission
	var offset = (page - 1) * limit

	r.DB = r.DB.Model(&data)
	r.DB = r.DB.Preload("Children").Preload("Children.Children")

	if err := r.DB.Where("parent_id IS NULL").Count(&count).Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return &data, count, nil
}

func (r repository) ReadByID(ID uint) (*entities.CrPermission, error) {
	var item entities.CrPermission

	r.DB = r.DB.Preload("Children").Preload("Children.Children")

	if err := r.DB.First(&item, ID).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r repository) Update(item *entities.CrPermission, payload *entities.CrPermission) error {
	if err := r.DB.Model(&item).Updates(payload).Error; err != nil {
		return err
	}

	return nil
}

func (r repository) Delete(ID uint) error {
	if err := r.DB.Unscoped().Delete(&entities.CrPermission{}, "id = ?", ID).Error; err != nil {
		return err
	}

	return nil
}
