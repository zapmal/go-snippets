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

	server := &http.Server{
		Addr:     config.Address,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	informationLog.Printf("Starting server on %s", config.Address)

	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
