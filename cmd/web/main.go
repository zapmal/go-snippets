package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"zapmal/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
)

type Config struct {
	Address         string
	StaticDirectory string
	DSN             string
	Secret          string
}

type Application struct {
	errorLog       *log.Logger
	informationLog *log.Logger
	session        *sessions.Session
	snippets       *mysql.SnippetModel
	users          *mysql.UserModel
	templateCache  map[string]*template.Template
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
	flag.StringVar(
		&config.Secret,
		"secret",
		getEnvVariable("SECRET_KEY"),
		"Secret Key",
	)
	flag.Parse()

	informationLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	database, err := openDatabaseConnection(config.DSN)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer database.Close()

	templateCache, err := newTemplateCache("./ui/html/")

	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(config.Secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &Application{
		errorLog:       errorLog,
		informationLog: informationLog,
		session:        session,
		snippets:       &mysql.SnippetModel{Database: database},
		users:          &mysql.UserModel{Database: database},
		templateCache:  templateCache,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr:         config.Address,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	informationLog.Printf("Starting server on %s", config.Address)

	// Needs to be generated with
	// go run go/path/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
