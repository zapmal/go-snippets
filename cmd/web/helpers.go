package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *Application) serverError(
	writer http.ResponseWriter,
	err error,
) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(
		writer,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func (app *Application) clientError(
	writer http.ResponseWriter,
	status int,
) {
	http.Error(writer, http.StatusText(status), status)
}

func (app *Application) notFound(writer http.ResponseWriter) {
	app.clientError(writer, http.StatusNotFound)
}

func (app *Application) addDefaultData(
	templateData *TemplateData,
	request *http.Request,
) *TemplateData {
	if templateData == nil {
		templateData = &TemplateData{}
	}

	templateData.CSRFToken = nosurf.Token(request)
	templateData.CurrentYear = time.Now().Year()
	templateData.Flash = app.session.PopString(request, "flashMessage")
	templateData.IsAuthenticated = app.isAuthenticated(request)

	return templateData
}

func (app *Application) render(
	writer http.ResponseWriter,
	request *http.Request,
	name string,
	templateData *TemplateData,
) {
	templateSet, ok := app.templateCache[name]

	if !ok {
		app.serverError(writer, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buffer := new(bytes.Buffer)

	err := templateSet.Execute(buffer, app.addDefaultData(templateData, request))

	if err != nil {
		app.serverError(writer, err)
		return
	}

	buffer.WriteTo(writer)
}

func (app *Application) isAuthenticated(request *http.Request) bool {
	isAuthenticated, ok := request.Context().Value(contextKeyIsAuthenticated).(bool)

	if !ok {
		return false
	}

	return isAuthenticated
}
