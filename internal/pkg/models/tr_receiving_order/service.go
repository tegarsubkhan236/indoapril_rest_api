package tr_receiving_order

import (
	"errors"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/ms_product_price"
	"example/internal/pkg/models/ms_stock"
	"example/internal/pkg/models/tr_back_order"
	"example/internal/pkg/models/tr_purchase_order"
	"example/internal/pkg/types/stock_status"
	types2 "example/internal/pkg/types/transaction_status"
	"example/internal/pkg/util/counter"
	"fmt"
	"strings"
	"time"
)

type Service interface {
	Create(req entities.TrReceivingOrderReq) (map[string]interface{}, error)

	createFromPO(req entities.TrReceivingOrderReq) (map[string]interface{}, error)
	createFromBO(req entities.TrReceivingOrderReq) (map[string]interface{}, error)

	validateProducts(po *entities.TrPurchaseOrder, bo *entities.TrBackOrder, req entities.TrReceivingOrderReq) (prices, quantities []int, quantityErrors error, remainingProducts []map[string]interface{})
	processRemainingProducts(po entities.TrPurchaseOrder, remainingProducts []map[string]interface{}) (poCode, boCode string, err error)
	incrementStockAndPrice(products []entities.TrReceivingOrderProductReq, userID uint) error
}

type service struct {
	stockRepo        ms_stock.Repository
	productPriceRepo ms_product_price.Repository
	poRepo           tr_purchase_order.Repository
	boRepo           tr_back_order.Repository
	roRepo           Repository
}

func NewService(
	a ms_stock.Repository,
	b ms_product_price.Repository,
	c tr_purchase_order.Repository,
	d tr_back_order.Repository,
	e Repository,
) Service {
	return &service{
		stockRepo:        a,
		productPriceRepo: b,
		poRepo:           c,
		boRepo:           d,
		roRepo:           e,
	}
}

func (s service) Create(req entities.TrReceivingOrderReq) (map[string]interface{}, error) {
	if req.PoCode != "" && req.BoCode == "" {
		return s.createFromPO(req)
	}

	if req.PoCode != "" && req.BoCode != "" {
		return s.createFromBO(req)
	}

	return nil, errors.New("params not enough")
}

func (s service) createFromPO(req entities.TrReceivingOrderReq) (map[string]interface{}, error) {
	tx := s.roRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			if tx.Error == nil {
				tx.Rollback()
			}
		} else if tx.Error == nil {
			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
			}
		}
	}()

	poItem, err := s.poRepo.ReadDetail(req.PoCode)
	if err != nil {
		return nil, err
	}

	if poItem.Status != types2.PENDING_ORDER {
		return nil, errors.New("po has been processed")
	}

	if len(poItem.TrPurchaseOrderProducts) != len(req.ReceivingOrderProducts) {
		return nil, errors.New("mismatched lengths")
	}

	prices, quantities, quantityError, remainingProducts := s.validateProducts(poItem, nil, req)
	if quantityError != nil {
		return nil, quantityError
	}

	req.PoID = poItem.ID
	req.Remarks = poItem.Remarks
	req.Amount, err = counter.CountAmount(poItem.Tax, poItem.Disc, prices, quantities)
	if err != nil {
		return nil, err
	}

	// Add Receiving Order
	payload := req.ToModel()
	if err = s.roRepo.CreateReceivingOrder(payload); err != nil {
		return nil, err
	}

	// Add Stock & Price
	if err = s.incrementStockAndPrice(req.ReceivingOrderProducts, req.UserID); err != nil {
		return nil, err
	}

	// Add Back Order
	if len(remainingProducts) > 0 {
		poCode, boCode, err := s.processRemainingProducts(*poItem, remainingProducts)
		if err != nil {
			return nil, err
		}

		if err = s.poRepo.Update(poItem.ID, entities.TrPurchaseOrder{Status: types2.PARTIAL_RECEIVE}); err != nil {
			return nil, err
		}

		result := map[string]interface{}{
			"message": fmt.Sprintf("purchase order with PO Code %s has been partially received and here is the new BO Code %s", poCode, boCode),
			"po_code": poCode,
			"bo_code": boCode,
		}

		return result, nil
	}

	// Update Status PO
	if err = s.poRepo.Update(poItem.ID, entities.TrPurchaseOrder{Status: types2.FULLY_RECEIVE}); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"message": fmt.Sprintf("purchase order with PO code %s has been fully received", poItem.PoCode),
		"po_code": poItem.PoCode,
		"bo_code": "",
	}

	return result, nil
}

func (s service) createFromBO(req entities.TrReceivingOrderReq) (map[string]interface{}, error) {
	tx := s.roRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			if tx.Error == nil {
				tx.Rollback()
			}
		} else if tx.Error == nil {
			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
			}
		}
	}()

	boItem, err := s.boRepo.ReadDetail(req.PoCode, req.BoCode)
	if err != nil {
		return nil, err
	}

	if boItem.Status != types2.PENDING_ORDER {
		return nil, errors.New("bo has been processed")
	}

	if len(boItem.TrBackOrderProducts) != len(req.ReceivingOrderProducts) {
		return nil, errors.New("mismatched lengths")
	}

	prices, quantities, quantityError, remainingProducts := s.validateProducts(nil, boItem, req)
	if quantityError != nil {
		return nil, quantityError
	}

	req.PoID = boItem.TrPurchaseOrderID
	req.BoID = &boItem.ID
	req.Remarks = boItem.Remarks
	req.Amount, err = counter.CountAmount(boItem.Tax, boItem.Disc, prices, quantities)
	if err != nil {
		return nil, err
	}

	// Add Receiving Order
	payload := req.ToModel()
	if err = s.roRepo.CreateReceivingOrder(payload); err != nil {
		return nil, err
	}

	// Add Stock & Price
	if err = s.incrementStockAndPrice(req.ReceivingOrderProducts, req.UserID); err != nil {
		return nil, err
	}

	// Add Back Order
	if len(remainingProducts) > 0 {
		poCode, boCode, err := s.processRemainingProducts(boItem.TrPurchaseOrder, remainingProducts)
		if err != nil {
			return nil, err
		}

		if err = s.boRepo.Update(boItem.ID, entities.TrBackOrder{Status: types2.PARTIAL_RECEIVE}); err != nil {
			return nil, err
		}

		result := map[string]interface{}{
			"message": fmt.Sprintf("back order with PO Code %s has been partially received and here is the new BO Code %s", poCode, boCode),
			"po_code": poCode,
			"bo_code": boCode,
		}

		return result, nil
	}

	// Update Status BO
	if err = s.boRepo.Update(boItem.ID, entities.TrBackOrder{Status: types2.FULLY_RECEIVE}); err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"message": fmt.Sprintf("back order with BO code %s has been fully received", boItem.BoCode),
		"po_code": boItem.TrPurchaseOrder.PoCode,
		"bo_code": boItem.BoCode,
	}

	return result, nil
}

func (s service) validateProducts(po *entities.TrPurchaseOrder, bo *entities.TrBackOrder, req entities.TrReceivingOrderReq) (prices, quantities []int, quantityErrors error, remainingProducts []map[string]interface{}) {
	var errStrings []string

	for i, val := range req.ReceivingOrderProducts {
		var productID uint
		var productName string
		var productQuantity, productPrice int

		if po != nil {
			productID = po.TrPurchaseOrderProducts[i].MsProduct.ID
			productName = po.TrPurchaseOrderProducts[i].MsProduct.Name
			productQuantity = po.TrPurchaseOrderProducts[i].Quantity
			productPrice = po.TrPurchaseOrderProducts[i].Price
			req.ReceivingOrderProducts[i].Price = po.TrPurchaseOrderProducts[i].Price
		}

		if bo != nil {
			productID = bo.TrBackOrderProducts[i].MsProduct.ID
			productName = bo.TrBackOrderProducts[i].MsProduct.Name
			productQuantity = bo.TrBackOrderProducts[i].Quantity
			productPrice = bo.TrBackOrderProducts[i].Price
			req.ReceivingOrderProducts[i].Price = bo.TrBackOrderProducts[i].Price
		}

		if val.Quantity > productQuantity {
			errorMessage := fmt.Errorf("the quantity of %s exceeds the required limit %d", productName, productQuantity)
			errStrings = append(errStrings, errorMessage.Error())
		}

		if val.Quantity < productQuantity {
			product := map[string]interface{}{
				"ProductID": productID,
				"Price":     productPrice,
				"Quantity":  productQuantity - val.Quantity,
			}
			remainingProducts = append(remainingProducts, product)
		}

		prices = append(prices, productPrice)
		quantities = append(quantities, val.Quantity)
	}

	if len(errStrings) > 0 {
		quantityErrors = fmt.Errorf(strings.Join(errStrings, "\n"))
	}

	return prices, quantities, quantityErrors, remainingProducts
}

func (s service) processRemainingProducts(po entities.TrPurchaseOrder, remainingProducts []map[string]interface{}) (poCode, boCode string, err error) {
	today := time.Now()
	lastBackOrderSequence, _ := s.boRepo.ReadLastItemSequence(today, po.ID)
	nextBackOrderSequence := lastBackOrderSequence + 1
	boCode = fmt.Sprintf("BO-%s%03d", today.Format("060102"), nextBackOrderSequence)

	boProductPrices := make([]int, 0)
	boProductQuantities := make([]int, 0)
	boProducts := make([]entities.TrBackOrderProduct, 0)

	for _, val := range remainingProducts {
		item := entities.TrBackOrderProduct{
			MsProductID: val["ProductID"].(uint),
			Price:       val["Price"].(int),
			Quantity:    val["Quantity"].(int),
		}
		boProducts = append(boProducts, item)
		boProductPrices = append(boProductPrices, val["Price"].(int))
		boProductQuantities = append(boProductQuantities, val["Quantity"].(int))
	}

	amount, err := counter.CountAmount(po.Tax, po.Disc, boProductPrices, boProductQuantities)
	if err != nil {
		return "", "", err
	}

	payload := entities.TrBackOrder{
		BoCode:              boCode,
		Disc:                po.Disc,
		Tax:                 po.Tax,
		Amount:              amount,
		Remarks:             po.Remarks,
		Status:              types2.PENDING_ORDER,
		TrPurchaseOrderID:   po.ID,
		MsSupplierID:        po.MsSupplierID,
		TrBackOrderProducts: boProducts,
	}

	if err = s.boRepo.Create(payload); err != nil {
		return "", "", err
	}

	return po.PoCode, boCode, nil
}

func (s service) incrementStockAndPrice(products []entities.TrReceivingOrderProductReq, userID uint) error {
	for _, val := range products {
		if err := s.stockRepo.CreateStock(stock_status.INCREMENT, userID, val.ProductID, val.Quantity); err != nil {
			return err
		}

		if err := s.productPriceRepo.CreateProductPrice(val.ProductID, val.Price); err != nil {
			return err
		}
	}

	return nil
}
