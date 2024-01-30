package tr_receiving_order

import (
	"example/internal/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	BeginTransaction() *gorm.DB
	CreateReceivingOrder(data entities.TrReceivingOrder) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) BeginTransaction() *gorm.DB {
	tx := r.DB.Begin()
	return tx
}

func (r repository) CreateReceivingOrder(data entities.TrReceivingOrder) error {
	err := r.DB.Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}
