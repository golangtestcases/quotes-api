package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangtestcases/quotes-api/handlers"
	"github.com/golangtestcases/quotes-api/repository"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func connectToDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/quotes?sslmode=disable"
	}

	var db *sql.DB
	var err error

	// Пытаемся подключиться в течение 30 секунд
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Attempt %d: DB connection error: %v", i+1, err)
			time.Sleep(3 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to PostgreSQL!")
			return db, nil
		}

		log.Printf("Attempt %d: DB ping failed: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to DB after %d attempts: %v", maxAttempts, err)
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS quotes (
			id SERIAL PRIMARY KEY,
			author TEXT NOT NULL,
			quote TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func main() {
	// Подключение к БД
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Fatal DB connection error: %v", err)
	}
	defer db.Close()

	// Создание таблиц
	if err := createTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Инициализация репозитория
	repo := repository.NewPostgresRepository(db)

	// Инициализация обработчиков
	quotesHandler := handlers.NewQuotesHandler(repo)

	// Настройка маршрутов
	r := mux.NewRouter()
	r.HandleFunc("/quotes", quotesHandler.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes", quotesHandler.AddQuote).Methods("POST")
	r.HandleFunc("/quotes/random", quotesHandler.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id}", quotesHandler.DeleteQuote).Methods("DELETE")

	// Настройка сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
