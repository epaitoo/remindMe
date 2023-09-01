package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/epaitoo/remindme/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	reminders     *models.ReminderModel
	templateCache map[string]*template.Template
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Mysql DSN
	dsn := flag.String("dsn", "web:password@/remindme?parseTime=true", "MySQL Data Source Name")
	flag.Parse()

	//loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create DB connection Pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err.Error())
	}

	defer db.Close()

	// initialize new templateCache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// New instance of the application struct
	app := &application{
		errorLog:  errorLog,
		infoLog:   infoLog,
		reminders: &models.ReminderModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
