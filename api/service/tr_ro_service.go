package service

import (
	"example/pkg"
	"example/pkg/model"
)

func CreateReceivingOrder(data model.TrReceivingOrder) error {
	var db = pkg.DB

	err := db.Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}
