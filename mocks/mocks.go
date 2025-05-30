package mocks

import (
	"github.com/golangtestcases/quotes-api/models"
	"github.com/golangtestcases/quotes-api/repository"
)

type MockRepository struct {
	repository.QuoteRepository
	AddFunc    func(quote models.Quote) (int, error)
	DeleteFunc func(id int) error
}

func (m *MockRepository) Add(quote models.Quote) (int, error) {
	return m.AddFunc(quote)
}

func (m *MockRepository) Delete(id int) error {
	return m.DeleteFunc(id)
}
