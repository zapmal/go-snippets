package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"zapmal/snippetbox/pkg/models"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("X-XSS-Protection", "1; mode=block")
			writer.Header().Set("X-Frame-Options", "deny")

			next.ServeHTTP(writer, request)
		},
	)
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			app.informationLog.Printf("%s - %s %s %s",
				request.RemoteAddr,
				request.Proto,
				request.Method,
				request.URL.RequestURI(),
			)

			next.ServeHTTP(writer, request)
		},
	)
}

func (app *Application) recoverFromPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			// This deferred function will always run in the event of a panic
			// because Go unwinds the Stack.
			defer func() {
				if err := recover(); err != nil {
					writer.Header().Set("Connection", "close")
					app.serverError(writer, fmt.Errorf("%s", err))
				}
			}()

			next.ServeHTTP(writer, request)
		},
	)
}

func (app *Application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			if !app.isAuthenticated(request) {
				http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
				return
			}

			writer.Header().Add("Cache-Control", "no-store")

			next.ServeHTTP(writer, request)
		},
	)
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			authenticationKeyExists := app.session.Exists(request, "authenticatedUserID")

			if !authenticationKeyExists {
				next.ServeHTTP(writer, request)
				return
			}

			user, err := app.users.Get(app.session.GetInt(request, "authenticatedUserID"))

			if errors.Is(err, models.ErrorRecordNotFound) || !user.Active {
				app.session.Remove(request, "authenticatedUserID")
				next.ServeHTTP(writer, request)
				return
			} else if err != nil {
				app.serverError(writer, err)
				return
			}

			ctx := context.WithValue(request.Context(), contextKeyIsAuthenticated, true)
			next.ServeHTTP(writer, request.WithContext(ctx))
		},
	)
}
