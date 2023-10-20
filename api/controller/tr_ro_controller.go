package controller

import (
	"example/api/service"
	"example/api/tool/converter"
	"example/api/tool/helper"
	"example/pkg/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func CreateReceivingOrder(c *fiber.Ctx) error {
	var req model.ReceivingOrderResponse
	if errParse := c.BodyParser(&req); errParse != nil {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "review your input", errParse.Error())
	}

	if req.PoCode == "" {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "po code cannot be empty", nil)
	}

	if req.UserID == 0 {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "user id cannot be empty", nil)
	}

	if req.PoCode != "" && req.BoCode == "" {
		return createReceivingOrderForPurchaseOrder(c, req)
	}

	if req.PoCode != "" && req.BoCode != "" {
		return createReceivingOrderForBackOrder(c, req)
	}

	return helper.ResponseHandler(c, fiber.StatusInternalServerError, "", nil)
}

func createReceivingOrderForPurchaseOrder(c *fiber.Ctx, req model.ReceivingOrderResponse) error {
	po, _, errPo := service.GetPurchaseOrder(0, 0, req.PoCode, "")
	if errPo != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, errPo.Error(), nil)
	}

	if po.Status != converter.PendingOrder {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "po has been processed", nil)
	}

	if len(po.TrPurchaseOrderProducts) != len(req.ReceivingOrderProducts) {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "mismatched lengths", nil)
	}

	prices, quantities, quantityError, remainingProducts := validateReceivingOrderProducts(po, model.TrBackOrder{}, req)

	if len(quantityError) > 0 {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "quantity errors", quantityError)
	}

	// Add Receiving Order
	req.PoID = po.ID
	req.Remarks = po.Remarks
	req.Amount, _ = helper.CountAmount(po.Tax, po.Disc, prices, quantities)
	roPayload := model.ConvertResponseToReceivingOrder(req)
	if errAddReceiving := service.CreateReceivingOrder(roPayload); errAddReceiving != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "receiving data failed to add", errAddReceiving.Error())
	}

	// Add Stock & Price
	if errAddStockAndPrice := incrementStockAndPrice(req.ReceivingOrderProducts, req.UserID); errAddStockAndPrice != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "stock and price failed to add", errAddStockAndPrice.Error())
	}

	// Add Back Order
	if len(remainingProducts) > 0 {
		err, poCode, boCode := processRemainingProducts(po, remainingProducts)
		if err != nil {
			return helper.ResponseHandler(c, fiber.StatusInternalServerError, "back order data failed to add", err.Error())
		}
		if errUpdatePO := service.UpdatePurchaseOrder(po.ID, model.TrPurchaseOrder{Status: converter.PartialReceive}); errUpdatePO != nil {
			return helper.ResponseHandler(c, fiber.StatusInternalServerError, "data failed to update", errUpdatePO.Error())
		}
		return helper.ResponseHandler(c, fiber.StatusOK, fmt.Sprintf("purchase order with code %s has been partially received", poCode), fiber.Map{
			"po_code": poCode,
			"bo_code": boCode,
		})
	}

	// Update Status PO
	if errUpdatePO := service.UpdatePurchaseOrder(po.ID, model.TrPurchaseOrder{Status: converter.FullyReceive}); errUpdatePO != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "data failed to update", errUpdatePO.Error())
	}

	return helper.ResponseHandler(c, fiber.StatusOK, fmt.Sprintf("purchase order with code %s has been fully received", req.PoCode), fiber.Map{
		"po_code": req.PoCode,
	})
}

func createReceivingOrderForBackOrder(c *fiber.Ctx, req model.ReceivingOrderResponse) error {
	_, bo, errBo := service.GetPurchaseOrder(0, 0, req.PoCode, req.BoCode)
	if errBo != nil {
		return helper.ResponseHandler(c, fiber.StatusNotFound, "back order not found", errBo.Error())
	}

	if bo.Status != converter.PendingOrder {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "bo has been processed", nil)
	}

	if len(bo.TrBackOrderProducts) != len(req.ReceivingOrderProducts) {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "mismatched lengths", nil)
	}

	prices, quantities, quantityError, remainingProducts := validateReceivingOrderProducts(model.TrPurchaseOrder{}, bo, req)
	if len(quantityError) > 0 {
		return helper.ResponseHandler(c, fiber.StatusBadRequest, "quantity errors", quantityError)
	}

	// Add Receiving Order
	req.PoID = bo.TrPurchaseOrderID
	req.BoID = bo.ID
	req.Remarks = bo.Remarks
	req.Amount, _ = helper.CountAmount(bo.Tax, bo.Disc, prices, quantities)
	roPayload := model.ConvertResponseToReceivingOrder(req)
	if errAddReceiving := service.CreateReceivingOrder(roPayload); errAddReceiving != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "receiving data failed to add", errAddReceiving.Error())
	}

	// Add Stock & Price
	if errAddStockAndPrice := incrementStockAndPrice(req.ReceivingOrderProducts, req.UserID); errAddStockAndPrice != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "stock and price failed to add", errAddStockAndPrice.Error())
	}

	// Add Next Back Order
	if len(remainingProducts) > 0 {
		err, poCode, boCode := processRemainingProducts(bo.TrPurchaseOrder, remainingProducts)
		if err != nil {
			return helper.ResponseHandler(c, fiber.StatusInternalServerError, "next back order data failed to add", bo.TrPurchaseOrder)
		}
		if errUpdateBO := service.UpdateBackOrder(bo.ID, model.TrBackOrder{Status: converter.PartialReceive}); errUpdateBO != nil {
			return helper.ResponseHandler(c, fiber.StatusInternalServerError, "data failed to update", errUpdateBO.Error())
		}
		return helper.ResponseHandler(c, fiber.StatusOK, fmt.Sprintf("purchase order with code %s has been partially received", poCode), fiber.Map{
			"po_code": poCode,
			"bo_code": boCode,
		})
	}

	// Update Status BO
	if errUpdateBO := service.UpdateBackOrder(bo.ID, model.TrBackOrder{Status: converter.FullyReceive}); errUpdateBO != nil {
		return helper.ResponseHandler(c, fiber.StatusInternalServerError, "data failed to update", errUpdateBO.Error())
	}

	return helper.ResponseHandler(c, fiber.StatusOK, fmt.Sprintf("purchase order with code %s has been fully received", req.BoCode), fiber.Map{
		"po_code": req.PoCode,
		"bo_code": req.BoCode,
	})
}

func validateReceivingOrderProducts(po model.TrPurchaseOrder, bo model.TrBackOrder, req model.ReceivingOrderResponse) ([]int, []int, []string, []map[string]interface{}) {
	var prices []int
	var quantities []int
	var quantityError []string
	var remainingProducts []map[string]interface{}
	var newRemainingProducts []map[string]interface{}

	for i, val := range req.ReceivingOrderProducts {
		var productID uint
		var productName string
		var productQuantity, productPrice int

		if po.ID != 0 {
			productID = po.TrPurchaseOrderProducts[i].MsProduct.ID
			productName = po.TrPurchaseOrderProducts[i].MsProduct.Name
			productQuantity = po.TrPurchaseOrderProducts[i].Quantity
			productPrice = po.TrPurchaseOrderProducts[i].Price
			req.ReceivingOrderProducts[i].Price = po.TrPurchaseOrderProducts[i].Price
		}
		if bo.ID != 0 {
			productID = bo.TrBackOrderProducts[i].MsProduct.ID
			productName = bo.TrBackOrderProducts[i].MsProduct.Name
			productQuantity = bo.TrBackOrderProducts[i].Quantity
			productPrice = bo.TrBackOrderProducts[i].Price
			req.ReceivingOrderProducts[i].Price = bo.TrBackOrderProducts[i].Price
		}

		if val.Quantity > productQuantity {
			errorMessage := fmt.Sprintf("the quantity of %s exceeds the required limit %d", productName, productQuantity)
			quantityError = append(quantityError, errorMessage)
		}
		if val.Quantity <= productQuantity {
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

	for _, remainingProduct := range remainingProducts {
		if quantity := remainingProduct["Quantity"].(int); quantity > 0 {
			newRemainingProducts = append(newRemainingProducts, remainingProduct)
		}
	}
	remainingProducts = newRemainingProducts

	return prices, quantities, quantityError, remainingProducts
}

func processRemainingProducts(po model.TrPurchaseOrder, remainingProducts []map[string]interface{}) (err error, poCode, boCode string) {
	today := time.Now()
	lastBackOrderSequence, _ := service.GetLastBoSequence(today, po.ID)
	nextBackOrderSequence := lastBackOrderSequence + 1
	boCode = fmt.Sprintf("BO-%s%03d", today.Format("060102"), nextBackOrderSequence)

	// Process back order products
	boProductPrices := make([]int, 0)
	boProductQuantities := make([]int, 0)
	boProducts := make([]model.TrBackOrderProduct, 0)

	for _, val := range remainingProducts {
		item := model.TrBackOrderProduct{
			MsProductID: val["ProductID"].(uint),
			Price:       val["Price"].(int),
			Quantity:    val["Quantity"].(int),
		}
		boProducts = append(boProducts, item)
		boProductPrices = append(boProductPrices, val["Price"].(int))
		boProductQuantities = append(boProductQuantities, val["Quantity"].(int))
	}

	amount, _ := helper.CountAmount(po.Tax, po.Disc, boProductPrices, boProductQuantities)
	payload := model.TrBackOrder{
		BoCode:              boCode,
		Disc:                po.Disc,
		Tax:                 po.Tax,
		Amount:              amount,
		Remarks:             po.Remarks,
		Status:              converter.PendingOrder,
		TrPurchaseOrderID:   po.ID,
		MsSupplierID:        po.MsSupplierID,
		TrBackOrderProducts: boProducts,
	}
	if errAddBackOrder := service.CreateBackOrder(payload); errAddBackOrder != nil {
		return errAddBackOrder, "", ""
	}

	return nil, po.PoCode, boCode
}

func incrementStockAndPrice(products []model.ReceivingOrderProductResponse, userID uint) error {
	for _, val := range products {
		errAddStock := service.CreateStock(converter.Increment, userID, val.ProductID, val.Quantity)
		if errAddStock != nil {
			return errAddStock
		}
		errAddPrice := service.CreateProductPrice(val.ProductID, val.Price)
		if errAddPrice != nil {
			return errAddPrice
		}
	}
	return nil
}
