package main

import (
	"ListItV3/pkg/app"
	"ListItV3/pkg/http/handlers"
	"ListItV3/pkg/repo"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=require password=%s port=%s search_path=%s", dbHost, dbUser, dbName, dbPassword, dbPort, "listitpoc")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Fatalf("Transaction rollback failed: %v", err)
		}
	}()
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS lists (id SERIAL PRIMARY KEY)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS items (id SERIAL PRIMARY KEY, name VARCHAR(255))")
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS list_items (id SERIAL PRIMARY KEY, list_id INTEGER REFERENCES lists(id), item_id INTEGER REFERENCES items(id))")
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_list_items_list_id ON list_items (list_id)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_list_items_item_id ON list_items (item_id)")
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	itemHandler := handlers.CreateHandler(app.NewItemSvc(repo.NewItemRepo(db)))
	listHandler := handlers.CreateListHandler(app.NewListSvc(repo.NewListRepo(db)))
	fmt.Println("Successfully connected! and repo created")
	r := chi.NewRouter()

	r.Route("/items", func(r chi.Router) {
		r.Post("/", itemHandler.CreateItem)
		r.Get("/{id}", itemHandler.GetItem)
		r.Delete("/{id}", itemHandler.DeleteItem)
		r.Get("/", itemHandler.GetAll)
	})
	r.Route("/list", func(r chi.Router) {
		r.Post("/", listHandler.CreateList)
		r.Get("/{id}", listHandler.GetList)
		r.Post("/{id}/addItem", listHandler.AddItemToList)
		//r.Get("/", itemHandler.GetAll)
	})

	http.ListenAndServe(":8008", r)
}
