package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
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
