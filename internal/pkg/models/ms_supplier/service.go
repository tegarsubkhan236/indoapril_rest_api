package ms_supplier

import (
	"example/internal/pkg/entities"
)

type Service interface {
	InsertSupplier(payload *[]entities.MsSupplier) (*[]entities.MsSupplier, error)
	FetchAllSupplier(page, limit int, ids []uint) (*[]entities.MsSupplier, int64, error)
	FetchDetailSupplier(ID uint) (*entities.MsSupplier, error)
	UpdateSupplier(ID uint, payload *entities.MsSupplier) (*entities.MsSupplier, error)
	DeleteSupplier(ID []uint) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s service) InsertSupplier(payload *[]entities.MsSupplier) (*[]entities.MsSupplier, error) {
	return s.repository.CreateSupplier(payload)
}

func (s service) FetchAllSupplier(page, limit int, ids []uint) (*[]entities.MsSupplier, int64, error) {
	return s.repository.ReadAllSupplier(page, limit, ids)
}

func (s service) FetchDetailSupplier(ID uint) (*entities.MsSupplier, error) {
	return s.repository.ReadSupplierById(ID)
}

func (s service) UpdateSupplier(ID uint, payload *entities.MsSupplier) (*entities.MsSupplier, error) {
	item, err := s.repository.ReadSupplierById(ID)
	if err != nil {
		return nil, err
	}

	return s.repository.UpdateSupplier(item, payload)
}

func (s service) DeleteSupplier(ID []uint) error {
	return s.repository.DestroySupplier(ID)
}
