package repo

import (
	"ListItV3/pkg/domain"
	"database/sql"
	"errors"
)

type ListRepo struct {
	db *sql.DB
}

func NewListRepo(db *sql.DB) *ListRepo {
	return &ListRepo{
		db: db,
	}
}

func (l ListRepo) Create(list *domain.List) (*domain.List, error) {
	if list.Items == nil {
		err := l.db.QueryRow("INSERT INTO lists default values RETURNING id").Scan(&list.ID)
		if err != nil {
			return nil, err
		} else {
			return list, nil
		}
	} else {
		tx, err := l.db.Begin()
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()

		err = tx.QueryRow("INSERT INTO lists default values RETURNING id").Scan(&list.ID)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, item := range list.Items {
			_, err := tx.Exec("INSERT INTO list_items (list_id, item_id) VALUES ($1, $2)", list.ID, item.ID)
			if err != nil {
				return nil, err
			}
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		err = tx.Commit()

		return list, err
	}
}

func (l ListRepo) Get(id int) (*domain.List, error) {
	list := &domain.List{ID: id}

	query := `
        SELECT i.id, i.name 
        FROM list_items li 
        JOIN items i ON li.item_id = i.id 
        WHERE li.list_id = $1
    `

	rows, err := l.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list.Items = []domain.Item{}
	for rows.Next() {
		item := domain.Item{}
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		list.Items = append(list.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If no rows were found, it means the list is empty or doesn't exist.
	return list, nil
}

func (l ListRepo) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}

func (l ListRepo) AddItem(id int, item domain.Item) error {
	// Check if the item ID is valid
	if item.ID == 0 {
		return errors.New("invalid item ID")
	}

	// Check if the item exists in the items table
	var exists bool
	err := l.db.QueryRow("SELECT exists (SELECT 1 FROM items WHERE id = $1)", item.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("item does not exist")
	}

	_, err = l.db.Exec("INSERT INTO list_items (list_id, item_id) VALUES ($1, $2)", id, item.ID)
	if err != nil {
		return err
	}

	return nil
}
