package tr_purchase_order

import (
	"errors"
	"example/internal/pkg/entities"
	"example/internal/pkg/models/tr_back_order"
	"example/internal/pkg/util/counter"
	"fmt"
	"time"
)

type Service interface {
	Insert(payload entities.TrPurchaseOrderReq) (*entities.TrPurchaseOrderResp, error)
	FetchAll(page, limit int) (*[]entities.TrPurchaseOrderResp, int64, error)
	FetchDetail(poCode, boCode *string) (*entities.TrPurchaseOrderResp, *entities.TrBackOrderResp, error)
}

type service struct {
	poRepo Repository
	boRepo tr_back_order.Repository
}

func NewService(poRepo Repository, boRepo tr_back_order.Repository) Service {
	return &service{
		poRepo: poRepo,
		boRepo: boRepo,
	}
}

func (s service) Insert(payload entities.TrPurchaseOrderReq) (*entities.TrPurchaseOrderResp, error) {
	today := time.Now()
	lastPoSequence, err := s.poRepo.ReadLastItemSequence(today)
	if err != nil {
		return nil, err
	}

	nextPoSequence := lastPoSequence + 1
	payload.PoCode = fmt.Sprintf("PO-%s%03d", today.Format("060102"), nextPoSequence)

	var prices []int
	var quantities []int
	for _, val := range payload.PurchaseOrderProducts {
		prices = append(prices, val.Price)
		quantities = append(quantities, val.Quantity)
	}
	payload.Amount, _ = counter.CountAmount(payload.Tax, payload.Disc, prices, quantities)

	item, err := s.poRepo.Create(payload.ToModel())
	if err != nil {
		return nil, err
	}

	result := item.ToResponse()
	return &result, nil
}

func (s service) FetchAll(page, limit int) (*[]entities.TrPurchaseOrderResp, int64, error) {
	data, count, err := s.poRepo.ReadAll(page, limit)
	if err != nil {
		return nil, 0, err
	}
	if data == nil {
		return nil, 0, nil
	}
	var results []entities.TrPurchaseOrderResp
	for _, item := range *data {
		results = append(results, item.ToResponse())
	}

	return &results, count, nil
}

func (s service) FetchDetail(poCode, boCode *string) (*entities.TrPurchaseOrderResp, *entities.TrBackOrderResp, error) {
	if poCode != nil && boCode == nil {
		poItem, err := s.poRepo.ReadDetail(*poCode)
		if err != nil {
			return nil, nil, err
		}
		result := poItem.ToResponse()

		return &result, nil, nil
	}
	if poCode != nil && boCode != nil {
		boItem, err := s.boRepo.ReadDetail(*poCode, *boCode)
		if err != nil {
			return nil, nil, err
		}
		result := boItem.ToResponse()

		return nil, &result, nil
	}

	return nil, nil, errors.New("params not enough")
}
