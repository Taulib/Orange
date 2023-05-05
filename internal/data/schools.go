package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/taulib/orange/internal/validator"
)

// School represents one row of data in our schools table
type School struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Level    string    `json:"level"`
	Contact  string    `json:"contact"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Website  string    `json:"website,omitempty"`
	Address  string    `json:"address"`
	Mode     []string  `json:"mode"`
	CreateAt time.Time `json:"-"`
	Version  int32     `json:"version"`
}

func ValidateSchool(v *validator.Validator, school *School) {
	v.Check(school.Name != "", "name", "must be provided")
	v.Check(len(school.Name) <= 200, "name", "Must not be more than 200 bytes long")

	v.Check(school.Level != "", "level", "must be provided")
	v.Check(len(school.Level) <= 200, "level", "Must not be more than 200 bytes long")

	v.Check(school.Contact != "", "contact", "must be provided")
	v.Check(len(school.Contact) <= 200, "contact", "Must not be more than 200 bytes long")

	v.Check(school.Phone != "", "phone", "must be provided")
	v.Check(validator.Matches(school.Phone, validator.PhoneRX), "phone", "Must be a valid phone number")

	v.Check(school.Email != "", "email", "must be provided")
	v.Check(validator.Matches(school.Email, validator.EmailRX), "email", "Must be a valid email address")

	v.Check(school.Website != "", "website", "must be provided")
	v.Check(validator.ValidWebsite(school.Website), "website", "Must be a valid url")

	v.Check(school.Address != "", "address", "must be provided")
	v.Check(len(school.Address) <= 500, "address", "Must not be more than 00 bytes long")

	v.Check(school.Mode != nil, "mode", "must be provided")
	v.Check(len(school.Mode) >= 1, "mode", "Must contain at least one mode")
	v.Check(len(school.Mode) <= 5, "mode", "Must contain at least five mode")
	v.Check(validator.Unique(school.Mode), "mode", "must not contain duplicate entries")
}

// Implement our models
// Define our Model
type SchoolModel struct {
	DB *sql.DB
}

// Insert a new school
func (m SchoolModel) Insert(school *School) error {
	//wrtie the sql code to insert
	query := `
			INSERT INTO schools(name, level, contact, phone, email, website, address, mode)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			RETURNING id, created_at , version		
	`
	// Collet the data fields into a slice
	args := []interface{}{school.Name, school.Level, school.Contact, school.Phone, school.Email,
		school.Website, school.Address, school.Mode}

	return m.DB.QueryRow(query, args...).Scan(&school.ID, &school.CreateAt, &school.Version)
}

// Get a new school
func (m SchoolModel) Get(id int64) (*School, error) {
	// Ensure that there is a valid ID
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Creating th Query
	query := `
	SELECT id, created_at, name, level, contact, phone, email, website, address, mode, version
	FROM schools
	WHERE id = $1
	`
	//declare a school variable to hold the returned data
	var school School
	// execute the query using QueryRow()
	err := m.DB.QueryRow(query, id).Scan(
		&school.ID,
		&school.CreateAt,
		&school.Name,
		&school.Level,
		&school.Contact,
		&school.Phone,
		&school.Email,
		&school.Website,
		&school.Address,
		pq.Array(&school.Mode),
		&school.Version,
	)
	// handle errors if any
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	//Sucess
	return &school, nil
}

// Update a school based on a identifier
// Optimistic Locking (version Number)

func (m SchoolModel) Update(school *School) error {
	// create a query
	query := `
	UPDATE schools
	SET name =$1, level = $2, contact =$3, phone = $4, email = $5, website = $6, 
					address =$7,mode = $8, version = version + 1
	WHERE id = $9
	AND version = $10
	RETURNING version
	`
	args := []interface{}{
		school.Name,
		school.Level,
		school.Contact,
		school.Phone,
		school.Email,
		school.Website,
		school.Address,
		pq.Array(school.Mode),
		school.ID,
		school.Version,
	}
	// check for edit conflicts
	err := m.DB.QueryRow(query, args...).Scan(&school.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete a new school based on an identifier
func (m SchoolModel) Delete(id int64) error {
	// Ensure that there is a valid ID
	if id < 1 {
		return ErrRecordNotFound
	}

	// Create the delete query
	query := `
	DELETE FROM schools
	WHERE id = $1
	`
	// Execute the query
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	// check how many rows where affected by the delete operation. We use the RowsAffected Method on the result variable

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrRecordNotFound
	}

	// check if no rows where affected
	if rowsAffected == 0 {
		return ErrRecordNotFound

	}
	return nil
}
