package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golangtestcases/quotes-api/models"
	"github.com/golangtestcases/quotes-api/repository"

	"github.com/gorilla/mux"
)

type QuotesHandler struct {
	repo repository.QuoteRepository
}

func NewQuotesHandler(repo repository.QuoteRepository) *QuotesHandler {
	return &QuotesHandler{repo: repo}
}

func (h *QuotesHandler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var quotes []models.Quote
	var err error

	if author != "" {
		quotes, err = h.repo.GetByAuthor(author)
	} else {
		quotes, err = h.repo.GetAll()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (h *QuotesHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.repo.GetRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

func (h *QuotesHandler) AddQuote(w http.ResponseWriter, r *http.Request) {
	var quote models.Quote
	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repo.Add(quote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	quote.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quote)
}

func (h *QuotesHandler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
