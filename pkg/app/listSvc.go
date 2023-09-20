package app

import "ListItV3/pkg/domain"

type ListSvc struct {
	repo domain.ListDB
}

func (l ListSvc) CreateList(list *domain.List) (*domain.List, error) {
	return l.repo.Create(list)
}

func (l ListSvc) GetList(id int) (*domain.List, error) {
	return l.repo.Get(id)
}

func (l ListSvc) DeleteList(id int) error {
	//TODO implement me
	panic("implement me")
}

func (l ListSvc) AddListItem(id int, item domain.Item) error {
	return l.repo.AddItem(id, item)
}

func NewListSvc(db domain.ListDB) domain.ListSvc {
	return &ListSvc{
		repo: db,
	}
}
