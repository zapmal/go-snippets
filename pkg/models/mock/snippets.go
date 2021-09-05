package mock

import (
	"time"

	"zapmal/snippetbox/pkg/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (model *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (model *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrorRecordNotFound
	}
}

func (model *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
