package repo

import (
	"ListItV3/pkg/domain"
	"database/sql"
	"fmt"
)

type ItemRepo struct {
	db *sql.DB
}

func NewItemRepo(db *sql.DB) *ItemRepo {
	return &ItemRepo{
		db: db,
	}
}

func (i ItemRepo) Create(item *domain.Item) (*domain.Item, error) {
	tx, err := i.db.Begin()
	if err != nil {
		return nil, err
	}
	err = tx.QueryRow("INSERT INTO items (name) VALUES ($1) RETURNING id", item.Name).Scan(&item.ID)

	if err != nil {
		tx.Rollback()

		return nil, err
	}
	if item.ID != 0 {
		tx.Commit()
		return item, nil
	} else {
		tx.Rollback()
		return nil, fmt.Errorf("error creating item more than one created, rolling back")
	}

}

func (i ItemRepo) Get(id int) (*domain.Item, error) {

	item := domain.Item{}
	err := i.db.QueryRow("SELECT id, name FROM items WHERE id = $1", id).Scan(&item.ID, &item.Name)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (i ItemRepo) Delete(id int) error {
	_, err := i.db.Exec("DELETE FROM items WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (i ItemRepo) GetAll() []*domain.Item {
	var items []*domain.Item
	rows, err := i.db.Query("SELECT id, name FROM items")
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		item := &domain.Item{}
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			return nil
		}
		items = append(items, item)
	}
	return items
}
