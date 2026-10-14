package repository

import "database/sql"

// Repository is the Postgres store for the example domain.
type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}
