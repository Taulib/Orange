package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("Edit Conflict")
)

// Create a wrapper for our models
type Models struct {
	Schools SchoolModel
}

// This function creates a Model instance
func NewModels(db *sql.DB) Models {
	return Models{
		Schools: SchoolModel{DB: db},
	}
}
