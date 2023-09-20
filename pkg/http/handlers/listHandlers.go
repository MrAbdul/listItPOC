package handlers

import (
	"ListItV3/pkg/domain"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type ListHandler struct {
	Svc domain.ListSvc
}

func CreateListHandler(svc domain.ListSvc) *ListHandler {
	return &ListHandler{
		Svc: svc,
	}
}

func (l *ListHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	var newList domain.List
	if err := json.NewDecoder(r.Body).Decode(&newList); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newList1, err := l.Svc.CreateList(&newList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newList1)

}

func (l *ListHandler) GetList(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	list, err := l.Svc.GetList(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (l *ListHandler) AddItemToList(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var newItem domain.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = l.Svc.AddListItem(id, newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
