package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Init initializes the database connection and creates the necessary tables.
func Init(dbURL string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 1; i <= 5; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("Waiting for the database to start...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	if err := createTable(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		finished BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
}
