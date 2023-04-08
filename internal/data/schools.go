package data

import (
	"time"

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
	v.Check(validator.Matches(school.Phone, validator.PhoneRx), "phone", "Must be a valid phone number")

	v.Check(school.Email != "", "email", "must be provided")
	v.Check(validator.Matches(school.Email, validator.EmailRx), "email", "Must be a valid email address")

	v.Check(school.Website != "", "website", "must be provided")
	v.Check(validator.ValidWebsite(school.Website), "website", "Must be a valid url")

	v.Check(school.Address != "", "address", "must be provided")
	v.Check(len(school.Address) <= 500, "address", "Must not be more than 00 bytes long")

	v.Check(school.Mode != nil, "mode", "must be provided")
	v.Check(len(school.Mode) >= 1, "mode", "Must contain at least one mode")
	v.Check(len(school.Mode) <= 5, "mode", "Must contain at least five mode")
	v.Check(validator.Unique(school.Mode), "mode", "must not contain duplicate entries")
}
