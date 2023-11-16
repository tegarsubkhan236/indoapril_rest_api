package tr_sales_order

import (
	"errors"
	"example/internal/pkg/entities"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Repository interface {
	CreateSalesOrder(data entities.TrSalesOrder) (*entities.TrSalesOrder, error)
	ReadSalesOrders(offset, limit int) (*[]entities.TrSalesOrder, int64, error)
	ReadSalesOrder(id uint, soCode ...string) (*entities.TrSalesOrder, error)
	ReadLastSoSequence(today time.Time) (int, error)
	UpdateSalesOrder(id uint, payload *entities.TrSalesOrder) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) CreateSalesOrder(data entities.TrSalesOrder) (*entities.TrSalesOrder, error) {
	err := r.DB.Create(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r repository) ReadSalesOrders(offset, limit int) (*[]entities.TrSalesOrder, int64, error) {
	var data []entities.TrSalesOrder
	var count int64

	r.DB = r.DB.Model(&data)

	r.DB = r.DB.Preload("TrSalesOrderProducts").Preload("TrSalesOrderProducts.MsProduct")

	if err := r.DB.Count(&count).Offset(offset).Limit(limit).Order("id desc").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return &data, count, nil
}

func (r repository) ReadSalesOrder(id uint, soCode ...string) (*entities.TrSalesOrder, error) {
	var item entities.TrSalesOrder

	r.DB = r.DB.Preload("TrSalesOrderProducts").Preload("TrSalesOrderProducts.MsProduct")

	if id > 0 {
		err := r.DB.First(&item, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no data found")
		}
	} else if len(soCode) > 0 {
		err := r.DB.Where("so_code = ?", soCode[0]).First(&item).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no data found")
		}
	} else {
		return nil, errors.New("either id or code should be provided")
	}

	return &item, nil
}

func (r repository) ReadLastSoSequence(today time.Time) (int, error) {
	var lastSo entities.TrSalesOrder
	var lastSoSequence int

	err := r.DB.Where("DATE(created_at) = ?", today.Format("2006-01-02")).Order("id desc").First(&lastSo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	_, _ = fmt.Sscanf(lastSo.SoCode, "SO-%s%03d", today.Format("060102"), &lastSoSequence)
	lastSoSuffix := lastSo.SoCode[len(lastSo.SoCode)-3:]
	lastSoSequence, _ = strconv.Atoi(lastSoSuffix)

	return lastSoSequence, nil
}

func (r repository) UpdateSalesOrder(id uint, payload *entities.TrSalesOrder) error {
	err := r.DB.Model(&entities.TrSalesOrder{}).Where("id = ?", id).Updates(payload).Error
	if err != nil {
		return err
	}

	return nil
}
