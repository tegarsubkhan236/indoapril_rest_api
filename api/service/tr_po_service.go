package service

import (
	"errors"
	"example/pkg"
	model2 "example/pkg/model"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func GetPurchaseOrders(page, limit int, filter model2.TrPurchaseOrder) ([]model2.TrPurchaseOrder, int64, error) {
	var db = pkg.DB
	var data []model2.TrPurchaseOrder
	var count int64
	var offset = (page - 1) * limit

	db = db.Model(&model2.TrPurchaseOrder{})

	if filter.PoCode != "" {
		db = db.Where("po_code = ?", filter.PoCode)
	}

	db = db.Preload("MsSupplier")
	db = db.Preload("TrPurchaseOrderProducts").Preload("TrPurchaseOrderProducts.MsProduct")
	db = db.Preload("TrBackOrders").Preload("TrBackOrders.TrBackOrderProducts").Preload("TrBackOrders.TrBackOrderProducts.MsProduct")

	result := db.Count(&count).Offset(offset).Limit(limit).Order("id desc").Find(&data)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return data, count, nil
}

func GetPurchaseOrder(poID, boID uint, poCode, boCode string) (model2.TrPurchaseOrder, model2.TrBackOrder, error) {
	var db = pkg.DB
	var purchaseOrderItem model2.TrPurchaseOrder
	var backOrderItem model2.TrBackOrder

	if poID != 0 {
		db = db.Preload("MsSupplier")
		db = db.Preload("TrPurchaseOrderProducts.MsProduct")
		db = db.Preload("TrBackOrders.TrBackOrderProducts.MsProduct")
		if err := db.First(&purchaseOrderItem, poID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return model2.TrPurchaseOrder{}, model2.TrBackOrder{}, errors.New("no data found")
			}
			return model2.TrPurchaseOrder{}, model2.TrBackOrder{}, err
		}
		return purchaseOrderItem, model2.TrBackOrder{}, nil
	}

	if boID != 0 {
		db = db.Preload("MsSupplier")
		db = db.Preload("TrBackOrderProducts")
		db = db.Preload("TrBackOrderProducts.MsProduct")
		if err := db.First(&backOrderItem, boID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return model2.TrPurchaseOrder{}, model2.TrBackOrder{}, errors.New("no data found")
			}
			return model2.TrPurchaseOrder{}, model2.TrBackOrder{}, err
		}
		return model2.TrPurchaseOrder{}, backOrderItem, nil
	}

	if poCode != "" && boCode == "" {
		db = db.Preload("MsSupplier")
		db = db.Preload("TrPurchaseOrderProducts.MsProduct")
		db = db.Preload("TrBackOrders.TrBackOrderProducts.MsProduct")
		if err := db.Where("po_code = ?", poCode).First(&purchaseOrderItem).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return model2.TrPurchaseOrder{}, model2.TrBackOrder{}, errors.New("no data found")
			}
			return model2.TrPurchaseOrder{}, model2.TrBackOrder{}, err
		}
		return purchaseOrderItem, model2.TrBackOrder{}, nil
	}

	if poCode != "" && boCode != "" {
		subQuery := db.Table("tr_purchase_orders").Select("id").Where("po_code = ?", poCode)
		db = db.Preload("MsSupplier")
		db = db.Preload("TrPurchaseOrder")
		db = db.Preload("TrBackOrderProducts")
		db = db.Preload("TrBackOrderProducts.MsProduct")
		if err := db.Where("bo_code = ? AND tr_purchase_order_id IN (?)", boCode, subQuery).First(&backOrderItem).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return model2.TrPurchaseOrder{}, model2.TrBackOrder{}, errors.New("no data found")
			}
			return model2.TrPurchaseOrder{}, model2.TrBackOrder{}, err
		}
		return model2.TrPurchaseOrder{}, backOrderItem, nil
	}

	return model2.TrPurchaseOrder{}, model2.TrBackOrder{}, errors.New("wrong function implementation")
}

func GetLastPoSequence(today time.Time) (int, error) {
	var lastPo model2.TrPurchaseOrder
	var lastPoSequence int
	var db = pkg.DB

	err := db.Where("DATE(created_at) = ?", today.Format("2006-01-02")).Order("id desc").First(&lastPo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	_, _ = fmt.Sscanf(lastPo.PoCode, "PO-%s%03d", today.Format("060102"), &lastPoSequence)
	lastPoSuffix := lastPo.PoCode[len(lastPo.PoCode)-3:]
	lastPoSequence, _ = strconv.Atoi(lastPoSuffix)
	return lastPoSequence, nil
}

func GetLastBoSequence(today time.Time, poID uint) (int, error) {
	var lastBackOrder model2.TrBackOrder
	var lastBackOrderSequence int
	var db = pkg.DB

	db = db.Where("tr_purchase_order_id = ?", poID)
	db = db.Where("DATE(created_at) = ?", today.Format("2006-01-02"))
	err := db.Order("id desc").First(&lastBackOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	_, _ = fmt.Sscanf(lastBackOrder.BoCode, "BO-%s%03d", today.Format("060102"), &lastBackOrderSequence)
	lastBackOrderSuffix := lastBackOrder.BoCode[len(lastBackOrder.BoCode)-3:]
	lastBackOrderSequence, _ = strconv.Atoi(lastBackOrderSuffix)
	return lastBackOrderSequence, nil
}

func CreatePurchaseOrder(data model2.TrPurchaseOrder) (model2.TrPurchaseOrder, error) {
	var db = pkg.DB

	err := db.Create(&data).Error
	if err != nil {
		return model2.TrPurchaseOrder{}, err
	}

	return data, nil
}

func CreateBackOrder(data model2.TrBackOrder) error {
	var db = pkg.DB

	err := db.Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdatePurchaseOrder(id uint, newData model2.TrPurchaseOrder) error {
	var db = pkg.DB

	err := db.Model(&model2.TrPurchaseOrder{}).Where("id = ?", id).Updates(newData).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateBackOrder(id uint, newData model2.TrBackOrder) error {
	var db = pkg.DB

	err := db.Model(&model2.TrBackOrder{}).Where("id = ?", id).Updates(newData).Error
	if err != nil {
		return err
	}

	return nil
}
