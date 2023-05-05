package validator

import (
	"net/url"
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	PhoneRX = regexp.MustCompile(`^\+?\(?[0-9]{3}\)?\s?-\s?[0-9]{3}\s?-\s?[0-9]{4}$`)
)

// We will create a map that will store validation errors
type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Methods that operate on our Validator type
// check if the map has any entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// add an entry to the map if the key does
// not already exist
func (v *Validator) AddError(key string, message string) {
	// check if the key is already in the map
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// check to see if an element can be found in a list
// of items
func In(element string, list ...string) bool {
	for i := range list {
		if element == list[i] {
			return true
		}
	}
	return false
}

// Matches() returns true if a string values matches
// a specific regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// ValidWebsite() checks if a string value is a valid URL
func ValidWebsite(website string) bool {
	_, err := url.ParseRequestURI(website)
	return err == nil
}

// Check() if we need to add an entry to the errors map
func (v *Validator) Check(ok bool, key string, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Unique() checks tha there are no repeating values in the slice
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
