package mysql

import (
	"reflect"
	"testing"
	"time"

	"zapmal/snippetbox/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name        string
		userID      int
		wantedUser  *models.User
		wantedError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantedUser: &models.User{
				ID:      1,
				Name:    "Alice Jones",
				Email:   "alice@example.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active:  true,
			},
			wantedError: nil,
		},
		{
			name:        "Zero ID",
			userID:      0,
			wantedUser:  nil,
			wantedError: models.ErrorRecordNotFound,
		},
		{
			name:        "Non-existent ID",
			userID:      2,
			wantedUser:  nil,
			wantedError: models.ErrorRecordNotFound,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			database, teardown := newTestDatabase(t)
			defer teardown()

			model := UserModel{database}

			user, err := model.Get(testCase.userID)

			if err != testCase.wantedError {
				t.Errorf("want %v; got %s", testCase.wantedError, err)
			}

			if !reflect.DeepEqual(user, testCase.wantedUser) {
				t.Errorf("want %v; got %v", testCase.wantedUser, user)
			}
		})
	}
}
