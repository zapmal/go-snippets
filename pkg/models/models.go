package models

import (
	"errors"
	"time"
)

var ErrorRecordNotFound = errors.New("models: No matching record found")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
