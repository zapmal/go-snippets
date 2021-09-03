package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverFromPanic, app.logRequest, secureHeaders)

	router := pat.New()

	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	router.Get("/", dynamicMiddleware.ThenFunc(app.home))
	router.Get("/snippet/create",
		dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm),
	)
	router.Post("/snippet/create",
		dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet),
	)
	router.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	router.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	router.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	router.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	router.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	router.Post("/user/logout",
		dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser),
	)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(router)
}
