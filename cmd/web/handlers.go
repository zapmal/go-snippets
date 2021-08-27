package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}

	writer.Write([]byte("Hello from SnippetBox."))
}

func showSnippet(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(writer, request)
		return
	}

	fmt.Fprintf(writer, "Display a specific snippet with ID %d", id)
}

func createSnippet(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}
	writer.Write([]byte("This is supposed to let you create a new snippet."))
}
