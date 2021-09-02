package mysql

import (
	"database/sql"

	"zapmal/snippetbox/pkg/models"
)

type UserModel struct {
	Database *sql.DB
}

func (userModel *UserModel) Insert(name, email, password string) error {
	return nil
}

func (userModel *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (userModel *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
