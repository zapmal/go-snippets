package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())

	defer server.Close()

	code, _, body := server.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnippet(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())

	defer server.Close()

	tests := []struct {
		name       string
		URLPath    string
		wantedCode int
		wantedBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			code, _, body := server.get(t, testCase.URLPath)

			if code != testCase.wantedCode {
				t.Errorf("want %d; got %d", testCase.wantedCode, code)
			}

			if !bytes.Contains(body, testCase.wantedBody) {
				t.Errorf("want body to caontain %q", testCase.wantedBody)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())

	defer server.Close()

	_, _, body := server.get(t, "/user/signup")
	csrfToken := extractCsrfToken(t, body)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantedCode   int
		wantedBody   []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field is too short (minimum is 10 characters")},
		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Address is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", testCase.userName)
			form.Add("email", testCase.userEmail)
			form.Add("password", testCase.userPassword)
			form.Add("csrf_token", testCase.csrfToken)

			code, _, body := server.postForm(t, "/user/signup", form)

			if code != testCase.wantedCode {
				t.Errorf("want %d; got %d", testCase.wantedCode, code)
			}

			if !bytes.Contains(body, testCase.wantedBody) {
				t.Errorf("want body %s to contain %q", body, testCase.wantedBody)
			}
		})
	}
}
