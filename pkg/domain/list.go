package domain

type List struct {
	ID    int    `json:"id"`
	Items []Item `json:"items"`
	//CreatedAt int64  `json:"createdAt"`
	//UpdatedAt int64  `json:"updatedAt"`
}

type ListSvc interface {
	CreateList(list *List) (*List, error)
	GetList(id int) (*List, error)
	DeleteList(id int) error
	AddListItem(id int, item Item) error
}

type ListDB interface {
	Create(list *List) (*List, error)
	Get(id int) (*List, error)
	Delete(id int) error
	AddItem(id int, item Item) error
}
