package storage

import (
	"database/sql"
	"fmt"
	"os"
)

// SQLiteStore handles database operations for the app
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLiteStore instance
func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

// InitializeSchema creates the database schema if it doesn't exist
func (s *SQLiteStore) InitializeSchema() error {
	// Check if tables already exist
	var tableName string
	err := s.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='categories' LIMIT 1").Scan(&tableName)
	if err == nil {
		// Table exists, schema is already initialized
		return nil
	} else if err != sql.ErrNoRows {
		// Some other error occurred
		return err
	}

	// Schema needs to be created
	schema, err := getSchemaSQL()
	if err != nil {
		return fmt.Errorf("failed to get schema SQL: %v", err)
	}

	// Execute the schema creation SQL
	_, err = s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}

	return nil
}

// getSchemaSQL reads the schema.sql file or returns a hardcoded schema
func getSchemaSQL() (string, error) {
	content, err := os.ReadFile("./schema.sql")
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// GetDB returns the database connection
func (s *SQLiteStore) GetDB() *sql.DB {
	return s.db
}
