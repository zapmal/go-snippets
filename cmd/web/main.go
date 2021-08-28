package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Address         string
	StaticDirectory string
}

type Application struct {
	errorLog       *log.Logger
	informationLog *log.Logger
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

	informationLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &Application{
		errorLog:       errorLog,
		informationLog: informationLog,
	}

	router := http.NewServeMux()
	router.HandleFunc("/", app.home)
	router.HandleFunc("/snippet", app.showSnippet)
	router.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("/static/", http.StripPrefix("/static", fileServer))

	server := &http.Server{
		Addr:     config.Address,
		ErrorLog: errorLog,
		Handler:  router,
	}
	informationLog.Printf("Starting server on %s", config.Address)

	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
