package domain

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ItemSvc interface {
	CreateItem(item *Item) (*Item, error)
	GetItem(id int) (*Item, error)
	DeleteItem(id int) error
	GetAll() []*Item
}

type ItemDB interface {
	Create(item *Item) (*Item, error)
	Get(id int) (*Item, error)
	Delete(id int) error
	GetAll() []*Item
}
