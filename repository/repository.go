package repository

import "github.com/golangtestcases/quotes-api/models"

type QuoteRepository interface {
	GetAll() ([]models.Quote, error)
	GetRandom() (models.Quote, error)
	GetByAuthor(author string) ([]models.Quote, error)
	Add(quote models.Quote) (int, error)
	Delete(id int) error
}
