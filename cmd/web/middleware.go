package main

import (
	"fmt"
	"net/http"
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
