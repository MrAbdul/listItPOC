package app

import (
	"ListItV3/pkg/domain"
)

type ItemSvc struct {
	repo domain.ItemDB
}

func (i ItemSvc) CreateItem(item *domain.Item) (*domain.Item, error) {

	return i.repo.Create(item)
}

func (i ItemSvc) GetItem(id int) (*domain.Item, error) {
	return i.repo.Get(id)
}

func (i ItemSvc) DeleteItem(id int) error {
	return i.repo.Delete(id)
}

func (i ItemSvc) GetAll() []*domain.Item {
	return i.repo.GetAll()
}

func NewItemSvc(db domain.ItemDB) domain.ItemSvc {
	return &ItemSvc{
		repo: db,
	}
}
