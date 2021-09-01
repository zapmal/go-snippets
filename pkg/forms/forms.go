package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (form *Form) Required(fields ...string) {
	for _, field := range fields {
		value := form.Get(field)
		if strings.TrimSpace(value) == "" {
			form.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (form *Form) MaxLength(field string, characterLimit int) {
	value := form.Get(field)

	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > characterLimit {
		form.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters", characterLimit))
	}
}

func (form *Form) AllowedValues(field string, options ...string) {
	value := form.Get(field)

	if value == "" {
		return
	}

	for _, option := range options {
		if value == option {
			return
		}
	}
	form.Errors.Add(field, "This field is invalid")
}

func (form *Form) Valid() bool {
	return len(form.Errors) == 0
}
