package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore wraps the standard sql.DB to provide specific repository methods.
// This follows the architectural pattern from the forum project.
type SQLiteStore struct {
	*sql.DB
}

// NewSQLiteStore initializes a new SQLite connection and applies the schema.
// It takes the dataSourceName (file path) and the path to the schema.sql file.
func NewSQLiteStore(dataSourceName, schemaFile string) (*SQLiteStore, error) {
	// Open the database connection
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Read the schema file from the filesystem
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		return nil, err
	}

	// Execute the schema to ensure all tables exist
	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db}, nil
}
