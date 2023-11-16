package ms_product

import (
	"example/internal/pkg/entities"
)

type Service interface {
	InsertBatchProduct(payload *[]entities.MsProductReq) (*[]entities.MsProductResp, error)
	FetchAllProduct(page, limit int, productName string, supplierID uint, productCategories []uint) (*[]entities.MsProductResp, int64, error)
	FetchDetailProduct(ID uint, basedOn string) (*entities.MsProductResp, error)
	UpdateProduct(ID uint, payload *entities.MsProductReq) (*entities.MsProductResp, error)
	DeleteProduct(ID []uint) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s service) InsertBatchProduct(payload *[]entities.MsProductReq) (*[]entities.MsProductResp, error) {
	var convPayload []entities.MsProduct
	for _, item := range *payload {
		convPayload = append(convPayload, item.ToModel())
	}

	data, err := s.repository.CreateBatchProduct(&convPayload)
	if err != nil {
		return nil, err
	}

	var results []entities.MsProductResp
	for _, item := range *data {
		results = append(results, item.ToResponse())
	}

	return &results, nil
}

func (s service) FetchAllProduct(page, limit int, productName string, supplierID uint, productCategories []uint) (*[]entities.MsProductResp, int64, error) {
	data, count, err := s.repository.ReadAllProduct(page, limit, productName, supplierID, productCategories)
	if err != nil {
		return nil, 0, err
	}

	if data == nil {
		return nil, 0, nil
	}

	var results []entities.MsProductResp
	for _, item := range *data {
		results = append(results, item.ToResponse())
	}

	return &results, count, nil
}

func (s service) FetchDetailProduct(ID uint, basedOn string) (*entities.MsProductResp, error) {
	item, err := s.repository.ReadProductById(ID, basedOn)
	if err != nil {
		return nil, err
	}

	result := item.ToResponse()
	return &result, nil
}

func (s service) UpdateProduct(ID uint, payload *entities.MsProductReq) (*entities.MsProductResp, error) {
	item, err := s.repository.ReadProductById(ID, "")
	if err != nil {
		return nil, err
	}

	convPayload := payload.ToModel()
	updatedItem, err := s.repository.UpdateProduct(item, &convPayload)
	if err != nil {
		return nil, err
	}

	result := updatedItem.ToResponse()
	return &result, nil
}

func (s service) DeleteProduct(ID []uint) error {
	err := s.repository.DestroyProduct(ID)
	if err != nil {
		return err
	}
	return nil
}
