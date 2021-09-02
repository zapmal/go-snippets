package models

import (
	"errors"
	"time"
)

var (
	ErrorRecordNotFound     = errors.New("models: No matching record found")
	ErrorInvalidCredentials = errors.New("models: Invalid credentials")
	ErrorDuplicateEmail     = errors.New("models: Duplicate email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password []byte
	Created  time.Time
	Active   bool
}
