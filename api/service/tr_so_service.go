package service

import (
	"errors"
	"example/pkg"
	"example/pkg/model"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func GetSalesOrders(offset, limit int, filter model.TrSalesOrder) ([]model.TrSalesOrder, int64, error) {
	var db = pkg.DB
	var data []model.TrSalesOrder
	var count int64

	db = db.Model(&model.TrSalesOrder{})

	if filter.SoCode != "" {
		db = db.Where("so_code = ?", filter.SoCode)
	}

	db = db.Preload("TrSalesOrderProducts").Preload("TrSalesOrderProducts.MsProduct")

	result := db.Count(&count).Offset(offset).Limit(limit).Order("id desc").Find(&data)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return data, count, nil
}

func GetSalesOrder(id uint, soCode ...string) (model.TrSalesOrder, error) {
	var db = pkg.DB
	var item model.TrSalesOrder

	db = db.Preload("TrSalesOrderProducts").Preload("TrSalesOrderProducts.MsProduct")

	if id > 0 {
		err := db.First(&item, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.TrSalesOrder{}, errors.New("no data found")
		}
	} else if len(soCode) > 0 {
		err := db.Where("so_code = ?", soCode[0]).First(&item).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.TrSalesOrder{}, errors.New("no data found")
		}
	} else {
		return model.TrSalesOrder{}, errors.New("either id or code should be provided")
	}

	return item, nil
}

func GetLastSoSequence(today time.Time) (int, error) {
	var lastSo model.TrSalesOrder
	var lastSoSequence int
	var db = pkg.DB

	err := db.Where("DATE(created_at) = ?", today.Format("2006-01-02")).Order("id desc").First(&lastSo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	_, _ = fmt.Sscanf(lastSo.SoCode, "SO-%s%03d", today.Format("060102"), &lastSoSequence)
	lastSoSuffix := lastSo.SoCode[len(lastSo.SoCode)-3:]
	lastSoSequence, _ = strconv.Atoi(lastSoSuffix)
	return lastSoSequence, nil
}

func CreateSalesOrder(data model.TrSalesOrder) (model.TrSalesOrder, error) {
	var db = pkg.DB

	err := db.Create(&data).Error
	if err != nil {
		return model.TrSalesOrder{}, err
	}

	return data, nil
}

func UpdateSalesOrder(id uint, newData model.TrSalesOrder) error {
	var db = pkg.DB

	err := db.Model(&model.TrSalesOrder{}).Where("id = ?", id).Updates(newData).Error
	if err != nil {
		return err
	}

	return nil
}
