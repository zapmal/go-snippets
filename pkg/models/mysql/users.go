package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"zapmal/snippetbox/pkg/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	Database *sql.DB
}

func (userModel *UserModel) Insert(name, email, password string) error {
	hashedPasswod, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	sqlStatement := `INSERT INTO users (name, email, password, created)
	VALUES(?,?,?, UTC_TIMESTAMP())
	`

	_, err = userModel.Database.Exec(sqlStatement, name, email, string(hashedPasswod))
	const DUPLICATE_ERROR_CODE = 1062

	if err != nil {
		var mySQLError *mysql.MySQLError

		if errors.As(err, &mySQLError) {
			if mySQLError.Number == DUPLICATE_ERROR_CODE &&
				strings.Contains(mySQLError.Message, "email_UNIQUE") {
				return models.ErrorDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (userModel *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	sqlStatement := "SELECT id, password FROM users WHERE email = ? and ACTIVE = TRUE"
	row := userModel.Database.QueryRow(sqlStatement, email)
	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrorInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrorInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (userModel *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
