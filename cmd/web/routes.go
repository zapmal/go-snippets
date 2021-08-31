package main

import "net/http"

func (app *Application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/", app.home)
	router.HandleFunc("/snippet", app.showSnippet)
	router.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("/static/", http.StripPrefix("/static", fileServer))

	return app.recoverFromPanic(app.logRequest(secureHeaders(router)))
}
