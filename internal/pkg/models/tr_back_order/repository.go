package tr_back_order

import (
	"errors"
	"example/internal/pkg/entities"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Repository interface {
	Create(data entities.TrBackOrder) error
	ReadDetail(poCode, boCode string) (*entities.TrBackOrder, error)
	ReadLastItemSequence(today time.Time, poID uint) (int, error)
	Update(id uint, newData entities.TrBackOrder) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) Create(data entities.TrBackOrder) error {
	if err := r.DB.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r repository) ReadDetail(poCode, boCode string) (*entities.TrBackOrder, error) {
	var backOrderItem entities.TrBackOrder
	subQuery := r.DB.Table("tr_purchase_orders").Select("id").Where("po_code = ?", poCode)

	r.DB = r.DB.Preload("MsSupplier")
	r.DB = r.DB.Preload("TrPurchaseOrder")
	r.DB = r.DB.Preload("TrBackOrderProducts")
	r.DB = r.DB.Preload("TrBackOrderProducts.MsProduct")

	if err := r.DB.Where("bo_code = ? AND tr_purchase_order_id IN (?)", boCode, subQuery).First(&backOrderItem).Error; err != nil {
		return nil, err
	}

	return &backOrderItem, nil
}

func (r repository) ReadLastItemSequence(today time.Time, poID uint) (int, error) {
	var lastBackOrder entities.TrBackOrder
	var lastBackOrderSequence int

	r.DB = r.DB.Where("tr_purchase_order_id = ?", poID)
	r.DB = r.DB.Where("DATE(created_at) = ?", today.Format("2006-01-02"))
	err := r.DB.Order("id desc").First(&lastBackOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	_, _ = fmt.Sscanf(lastBackOrder.BoCode, "BO-%s%03d", today.Format("060102"), &lastBackOrderSequence)
	lastBackOrderSuffix := lastBackOrder.BoCode[len(lastBackOrder.BoCode)-3:]
	lastBackOrderSequence, _ = strconv.Atoi(lastBackOrderSuffix)
	return lastBackOrderSequence, nil
}

func (r repository) Update(id uint, newData entities.TrBackOrder) error {
	err := r.DB.Model(&entities.TrBackOrder{}).Where("id = ?", id).Updates(newData).Error
	if err != nil {
		return err
	}

	return nil
}
