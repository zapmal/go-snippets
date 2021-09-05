package main

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"zapmal/snippetbox/pkg/models/mock"

	"github.com/golangcollege/sessions"
)

var csrfTokenRegexp = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)"/>`)

func extractCsrfToken(t *testing.T, body []byte) string {
	matches := csrfTokenRegexp.FindSubmatch(body)

	if len(matches) < 2 {
		t.Fatal("No CSRF Token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}

func newTestApplication(t *testing.T) *Application {
	templateCache, err := newTemplateCache("./../../ui/html/")

	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return &Application{
		errorLog:       log.New(ioutil.Discard, "", 0),
		informationLog: log.New(ioutil.Discard, "", 0),
		session:        session,
		snippets:       &mock.SnippetModel{},
		templateCache:  templateCache,
		users:          &mock.UserModel{},
	}
}

type TestServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *TestServer {
	server := httptest.NewTLSServer(h)

	cookieJar, err := cookiejar.New(nil)

	if err != nil {
		t.Fatal(err)
	}

	server.Client().Jar = cookieJar

	server.Client().CheckRedirect = func(request *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestServer{server}
}

func (server *TestServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	response, err := server.Client().Get(server.URL + urlPath)

	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}

	return response.StatusCode, response.Header, body
}

func (server *TestServer) postForm(
	t *testing.T,
	URLPath string,
	form url.Values,
) (int, http.Header, []byte) {
	response, err := server.Client().PostForm(server.URL+URLPath, form)

	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}

	return response.StatusCode, response.Header, body
}
