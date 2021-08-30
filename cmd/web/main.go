package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Config struct {
	Address         string
	StaticDirectory string
	DSN             string
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
	flag.StringVar(
		&config.DSN,
		"dsn",
		getEnvVariable("DATABASE_DSN"),
		"MySQL data source name",
	)
	flag.Parse()

	informationLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	connection, err := openDatabaseConnection(config.DSN)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer connection.Close()

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

	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func openDatabaseConnection(DSN string) (*sql.DB, error) {
	connection, err := sql.Open("mysql", DSN)

	if err != nil {
		return nil, err
	}

	if err = connection.Ping(); err != nil {
		return nil, err
	}

	return connection, nil
}

func getEnvVariable(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Couldn't load .env file.")
	}

	return os.Getenv(key)
}
