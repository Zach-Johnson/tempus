package storage

import (
	"database/sql"
)

// SQLiteStore handles database operations for the app
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLiteStore instance
func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

// GetDB returns the database connection
func (s *SQLiteStore) GetDB() *sql.DB {
	return s.db
}
