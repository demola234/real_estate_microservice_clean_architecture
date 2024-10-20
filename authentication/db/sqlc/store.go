package db

import (
	"database/sql"
)

type Store interface {
	// Methods from the Querier interface (generated by sqlc)
	Querier
}

// SQLStore implements the Store interface and provides transaction support.
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new SQLStore instance with the provided DB and returns it as Store.
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db), // Assumes New is an sqlc-generated constructor for Queries
	}
}
