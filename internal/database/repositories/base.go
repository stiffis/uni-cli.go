package repositories

import (
	"database/sql"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	db *sql.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// DB returns the underlying database connection
func (r *BaseRepository) DB() *sql.DB {
	return r.db
}

// BeginTx starts a new transaction
func (r *BaseRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}
