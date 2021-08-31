package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"zapmal/snippetbox/pkg/models"
)

func (app *Application) home(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.URL.Path != "/" {
		app.notFound(writer)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.render(writer, request, "home.page.tmpl", &TemplateData{
		Snippets: snippets,
	})
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

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrorRecordNotFound) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}

		return
	}

	app.render(writer, request, "show.page.tmpl", &TemplateData{
		Snippet: snippet,
	})
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

	title := "0 Snail"
	content := "Who know what did the snail do? Probably nothing but whatever"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(writer, err)
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
