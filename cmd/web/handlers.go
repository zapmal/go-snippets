package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"zapmal/snippetbox/pkg/forms"
	"zapmal/snippetbox/pkg/models"
)

func (app *Application) home(
	writer http.ResponseWriter,
	request *http.Request,
) {
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
	id, err := strconv.Atoi(request.URL.Query().Get(":id"))

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

func (app *Application) createSnippetForm(
	writer http.ResponseWriter,
	request *http.Request,
) {
	app.render(writer, request, "create.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

func (app *Application) createSnippet(
	writer http.ResponseWriter,
	request *http.Request,
) {
	err := request.ParseForm()

	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.AllowedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(writer, request, "create.page.tmpl", &TemplateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(
		form.Get("title"),
		form.Get("content"),
		form.Get("expires"),
	)

	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.session.Put(request, "flashMessage", "Snippet created successfully")

	http.Redirect(writer, request, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *Application) signupUserForm(
	writer http.ResponseWriter,
	request *http.Request,
) {
	app.render(writer, request, "signup.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

func (app *Application) signupUser(
	writer http.ResponseWriter,
	request *http.Request,
) {
	err := request.ParseForm()

	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.ValidateEmail()
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(writer, request, "signup.page.tmpl", &TemplateData{
			Form: form,
		})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))

	if err != nil {
		if errors.Is(err, models.ErrorDuplicateEmail) {
			form.Errors.Add("email", "Email address is already in use")
			app.render(writer, request, "signup.page.tmpl", &TemplateData{Form: form})
		} else {
			app.serverError(writer, err)
		}

		return
	}

	app.session.Put(request, "flashMessage", "Your signup was successful. You can login now.")
	http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
}

func (app *Application) loginUserForm(
	writer http.ResponseWriter,
	request *http.Request,
) {
	app.render(writer, request, "login.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

func (app *Application) loginUser(
	writer http.ResponseWriter,
	request *http.Request,
) {
	err := request.ParseForm()

	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := forms.New(request.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))

	if err != nil {
		if errors.Is(err, models.ErrorInvalidCredentials) {
			form.Errors.Add("generic", "Email and/or password is incorrect")
			app.render(writer, request, "login.page.tmpl", &TemplateData{Form: form})
		} else {
			app.serverError(writer, err)
		}

		return
	}

	app.session.Put(request, "authenticatedUserID", id)
	http.Redirect(writer, request, "/snippet/create", http.StatusSeeOther)
}

func (app *Application) logoutUser(
	writer http.ResponseWriter,
	request *http.Request,
) {
	app.session.Remove(request, "authenticatedUserID")

	app.session.Put(request, "flashMessage", "You've been logged out successfully!")
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
