package main

import (
	"flag"
	"log"
	"net/http"
)

type Config struct {
	Address         string
	StaticDirectory string
}

func main() {
	config := new(Config)
	flag.StringVar(&config.Address, "addr", ":4000", "HTTP Network Address")
	flag.StringVar(
		&config.StaticDirectory,
		"static-dir",
		"./ui/static",
		"Path to static assets",
	)
	flag.Parse()

	router := http.NewServeMux()
	router.HandleFunc("/", home)
	router.HandleFunc("/snippet", showSnippet)
	router.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting server on %s", config.Address)
	err := http.ListenAndServe(config.Address, router)

	log.Fatal(err)
}
