package mysql

import (
	"database/sql"
	"io/ioutil"
	"testing"
)

func newTestDatabase(t *testing.T) (*sql.DB, func()) {
	database, err := sql.Open("mysql", "test_manuel:test_password@/test_snippetbox?parseTime=true&multiStatements=true")

	if err != nil {
		t.Fatal(err)
	}

	script, err := ioutil.ReadFile("./testdata/setup.sql")

	if err != nil {
		t.Fatal(err)
	}

	_, err = database.Exec(string(script))

	if err != nil {
		t.Fatal(err)
	}

	return database, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")

		if err != nil {
			t.Fatal(err)
		}

		_, err = database.Exec(string(script))

		if err != nil {
			t.Fatal(err)
		}

		database.Close()
	}
}
