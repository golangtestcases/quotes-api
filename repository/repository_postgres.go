package repository

import (
	"database/sql"

	"github.com/golangtestcases/quotes-api/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAll() ([]models.Quote, error) {
	rows, err := r.db.Query("SELECT id, author, quote FROM quotes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quotes []models.Quote
	for rows.Next() {
		var q models.Quote
		if err := rows.Scan(&q.ID, &q.Author, &q.Quote); err != nil {
			return nil, err
		}
		quotes = append(quotes, q)
	}
	return quotes, nil
}

func (r *PostgresRepository) GetRandom() (models.Quote, error) {
	var q models.Quote
	err := r.db.QueryRow("SELECT id, author, quote FROM quotes ORDER BY RANDOM() LIMIT 1").Scan(&q.ID, &q.Author, &q.Quote)
	return q, err
}

func (r *PostgresRepository) GetByAuthor(author string) ([]models.Quote, error) {
	rows, err := r.db.Query("SELECT id, author, quote FROM quotes WHERE author = $1", author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quotes []models.Quote
	for rows.Next() {
		var q models.Quote
		if err := rows.Scan(&q.ID, &q.Author, &q.Quote); err != nil {
			return nil, err
		}
		quotes = append(quotes, q)
	}
	return quotes, nil
}

func (r *PostgresRepository) Add(quote models.Quote) (int, error) {
	var id int
	err := r.db.QueryRow("INSERT INTO quotes (author, quote) VALUES ($1, $2) RETURNING id", quote.Author, quote.Quote).Scan(&id)
	return id, err
}

func (r *PostgresRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM quotes WHERE id = $1", id)
	return err
}
