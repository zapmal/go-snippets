package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverFromPanic, app.logRequest, secureHeaders)

	router := pat.New()

	router.Get("/", http.HandlerFunc(app.home))
	router.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	router.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	router.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(router)
}
