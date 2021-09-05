package forms

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRegexp = regexp.MustCompile("/^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/")

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

func (form *Form) MinLength(field string, minLength int) {
	value := form.Get(field)

	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) < minLength {
		form.Errors.Add(field,
			fmt.Sprintf("This field is too short (minimum is %d characters)", minLength),
		)
	}
}

func (form *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := form.Get(field)

	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		form.Errors.Add(field, "This field is invalid")
	}
}

func (form *Form) ValidateEmail() {
	value := form.Get("email")

	if value == "" {
		return
	}

	_, err := mail.ParseAddress(value)

	if err != nil {
		form.Errors.Add("email", "This field is invalid")
	}
}
