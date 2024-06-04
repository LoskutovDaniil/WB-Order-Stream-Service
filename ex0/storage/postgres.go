package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgres() (db *sql.DB, err error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", "localhost", "5432", "postgres", "ex0", "disable")

	if db, err = sql.Open("postgres", url); err != nil {
		return nil, err 
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
