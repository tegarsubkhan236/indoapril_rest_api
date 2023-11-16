package ms_product_category

import "example/internal/pkg/entities"

type Service interface {
	InsertProductCategory(payload *entities.MsProductCategory) error
	FetchAllProductCategory(page, limit int) (*[]entities.MsProductCategoryResp, int64, error)
	FetchDetailProductCategory(ID uint) (*entities.MsProductCategoryResp, error)
	UpdateProductCategory(ID uint, payload *entities.MsProductCategory) error
	DeleteProductCategory(ID uint) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) InsertProductCategory(payload *entities.MsProductCategory) error {
	return s.repository.CreateProductCategory(payload)
}

func (s service) FetchAllProductCategory(page, limit int) (*[]entities.MsProductCategoryResp, int64, error) {
	productCategories, count, err := s.repository.ReadAllProductCategory(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var results []entities.MsProductCategoryResp
	for _, item := range *productCategories {
		results = append(results, item.ToResponse())
	}

	return &results, count, nil
}

func (s service) FetchDetailProductCategory(ID uint) (*entities.MsProductCategoryResp, error) {
	productCategory, err := s.repository.ReadProductCategoryById(ID)
	if err != nil {
		return nil, err
	}
	result := productCategory.ToResponse()
	return &result, nil
}

func (s service) UpdateProductCategory(ID uint, payload *entities.MsProductCategory) error {
	item, err := s.repository.ReadProductCategoryById(ID)
	if err != nil {
		return err
	}
	return s.repository.UpdateProductCategory(item, payload)
}

func (s service) DeleteProductCategory(ID uint) error {
	return s.repository.DestroyProductCategory(ID)
}
