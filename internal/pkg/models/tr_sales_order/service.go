package tr_sales_order

import (
	"example/internal/pkg/entities"
	"example/internal/pkg/models/ms_product_price"
	"example/internal/pkg/models/ms_stock"
	"example/internal/pkg/types/stock_status"
	"example/internal/pkg/util/counter"
	"fmt"
	"strings"
	"time"
)

type Service interface {
	Insert(req entities.TrSalesOrderReq) (*entities.TrSalesOrderResp, error)
	FetchAll(page, limit int) (*[]entities.TrSalesOrderResp, int64, error)
	FetchDetail(id uint) (*entities.TrSalesOrderResp, error)

	decrementStock(products []entities.TrSalesOrderProductReq, userID uint) error
}

type service struct {
	soRepo           Repository
	stockRepo        ms_stock.Repository
	productPriceRepo ms_product_price.Repository
}

func NewService(r Repository, s ms_stock.Repository, p ms_product_price.Repository) Service {
	return &service{
		soRepo:           r,
		stockRepo:        s,
		productPriceRepo: p,
	}
}

func (s service) Insert(req entities.TrSalesOrderReq) (*entities.TrSalesOrderResp, error) {
	today := time.Now()
	lastSoSequence, _ := s.soRepo.ReadLastSoSequence(today)
	nextSoSequence := lastSoSequence + 1
	soCode := fmt.Sprintf("SO-%s%03d", today.Format("060102"), nextSoSequence)

	req.SoCode = soCode

	var prices []int
	var quantities []int
	for _, val := range req.SalesOrderProducts {
		productPrice, _ := s.productPriceRepo.ReadLastProductPrice(val.ProductID)
		prices = append(prices, productPrice.SellPrice)
		quantities = append(quantities, val.Quantity)
	}

	amount, err := counter.CountAmount(req.Tax, req.Disc, prices, quantities)
	if err != nil {
		return nil, err
	}

	req.Amount = amount
	payload := req.ToModel()
	item, err := s.soRepo.CreateSalesOrder(payload)
	if err != nil {
		return nil, err
	}

	if err = s.decrementStock(req.SalesOrderProducts, req.UserID); err != nil {
		return nil, err
	}

	result := item.ToResponse()
	return &result, nil
}

func (s service) FetchAll(page, limit int) (*[]entities.TrSalesOrderResp, int64, error) {
	data, count, err := s.soRepo.ReadSalesOrders(page, limit)
	if err != nil {
		return nil, 0, err
	}
	if data == nil {
		return nil, 0, nil
	}
	var results []entities.TrSalesOrderResp
	for _, item := range *data {
		results = append(results, item.ToResponse())
	}

	return &results, count, nil
}

func (s service) FetchDetail(id uint) (*entities.TrSalesOrderResp, error) {
	item, err := s.soRepo.ReadSalesOrder(id)
	if err != nil {
		return nil, err
	}
	result := item.ToResponse()

	return &result, nil
}

func (s service) decrementStock(products []entities.TrSalesOrderProductReq, userID uint) error {
	var errString []string

	for _, val := range products {
		if err := s.stockRepo.CreateStock(stock_status.DECREMENT, userID, val.ProductID, val.Quantity); err != nil {
			errString = append(errString, err.Error())
		}
	}

	if len(errString) > 0 {
		return fmt.Errorf(strings.Join(errString, "\n"))
	}

	return nil
}
