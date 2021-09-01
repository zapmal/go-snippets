package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverFromPanic, app.logRequest, secureHeaders)

	router := pat.New()

	dynamicMiddleware := alice.New(app.session.Enable)

	router.Get("/", dynamicMiddleware.ThenFunc(app.home))
	router.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	router.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	router.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(router)
}
