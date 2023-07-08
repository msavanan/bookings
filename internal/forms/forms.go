package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes the form structure
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Has(field string) bool {
	x := f.Get(field)

	if x == "" {
		return false
	}

	return true

}

// Validation for the required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Validation for the minimum length of field
func (f *Form) MinimumLength(field string, length int) {
	value := f.Get(field)

	if len(value) < length {
		f.Errors.Add(field, fmt.Sprintf("This field should have a minimum length of %d charachters long", length))
	}
}

//Validate Email

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
