package mysql

import (
	"database/sql"
	"errors"

	"zapmal/snippetbox/pkg/models"
)

type SnippetModel struct {
	Database *sql.DB
}

func (model *SnippetModel) Insert(
	title,
	content,
	expires string,
) (int, error) {
	sqlStatement := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := model.Database.Exec(sqlStatement, title, content, expires)

	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (model *SnippetModel) Get(id int) (*models.Snippet, error) {
	sqlStatement := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	snippet := &models.Snippet{}
	err := model.Database.QueryRow(sqlStatement, id).
		Scan(
			&snippet.ID,
			&snippet.Title,
			&snippet.Content,
			&snippet.Created,
			&snippet.Expires,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrorRecordNotFound
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

func (model *SnippetModel) Latest() ([]*models.Snippet, error) {
	sqlStatement := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10
	`

	rows, err := model.Database.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		snippet := &models.Snippet{}
		err = rows.Scan(
			&snippet.ID,
			&snippet.Title,
			&snippet.Content,
			&snippet.Created,
			&snippet.Expires,
		)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
