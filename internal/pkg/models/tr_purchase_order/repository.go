package tr_purchase_order

import (
	"errors"
	"example/internal/pkg/entities"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Repository interface {
	Create(data entities.TrPurchaseOrder) (*entities.TrPurchaseOrder, error)
	ReadAll(page, limit int) (*[]entities.TrPurchaseOrder, int64, error)
	ReadDetail(poCode string) (*entities.TrPurchaseOrder, error)
	ReadLastItemSequence(today time.Time) (int, error)
	Update(id uint, newData entities.TrPurchaseOrder) error
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) ReadAll(page, limit int) (*[]entities.TrPurchaseOrder, int64, error) {
	var data []entities.TrPurchaseOrder
	var count int64
	var offset = (page - 1) * limit

	r.DB = r.DB.Model(&data)

	r.DB = r.DB.Preload("MsSupplier")
	r.DB = r.DB.Preload("TrPurchaseOrderProducts").
		Preload("TrPurchaseOrderProducts.MsProduct")
	r.DB = r.DB.Preload("TrBackOrders").
		Preload("TrBackOrders.TrBackOrderProducts").
		Preload("TrBackOrders.TrBackOrderProducts.MsProduct")

	if err := r.DB.Count(&count).Offset(offset).Limit(limit).Order("id desc").Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return &data, count, nil
}

func (r repository) ReadDetail(poCode string) (*entities.TrPurchaseOrder, error) {
	var purchaseOrderItem entities.TrPurchaseOrder

	r.DB = r.DB.Preload("MsSupplier")
	r.DB = r.DB.Preload("TrPurchaseOrderProducts.MsProduct")
	r.DB = r.DB.Preload("TrBackOrders.TrBackOrderProducts.MsProduct")

	if err := r.DB.Where("po_code = ?", poCode).First(&purchaseOrderItem).Error; err != nil {
		return nil, err
	}

	return &purchaseOrderItem, nil
}

func (r repository) ReadLastItemSequence(today time.Time) (int, error) {
	var lastPo entities.TrPurchaseOrder
	var lastPoSequence int

	err := r.DB.Where("DATE(created_at) = ?", today.Format("2006-01-02")).Order("id desc").First(&lastPo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	_, _ = fmt.Sscanf(lastPo.PoCode, "PO-%s%03d", today.Format("060102"), &lastPoSequence)
	lastPoSuffix := lastPo.PoCode[len(lastPo.PoCode)-3:]
	lastPoSequence, _ = strconv.Atoi(lastPoSuffix)

	return lastPoSequence, nil
}

func (r repository) Create(data entities.TrPurchaseOrder) (*entities.TrPurchaseOrder, error) {
	if err := r.DB.Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r repository) Update(id uint, item entities.TrPurchaseOrder) error {
	if err := r.DB.Model(&item).Where("id = ?", id).Updates(item).Error; err != nil {
		return err
	}

	return nil
}
