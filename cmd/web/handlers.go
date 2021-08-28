package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *Application) home(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.URL.Path != "/" {
		app.notFound(writer)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	templateSet, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(writer, err)
		return
	}

	err = templateSet.Execute(writer, nil)

	if err != nil {
		app.serverError(writer, err)
	}
}

func (app *Application) showSnippet(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	fmt.Fprintf(writer, "Display a specific snippet with ID %d", id)
}

func (app *Application) createSnippet(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)

		return
	}
	writer.Write([]byte("This is supposed to let you create a new snippet."))
}
