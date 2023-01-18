package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golangify.com/snippetbox/pkg/models/sqlite3"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *sqlite3.SnippetModel
}

const dsn = "store.db"

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &sqlite3.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\nhttp://localhost:%s/", *addr, *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err := createDB(db); err != nil {
		return nil, err
	}

	return db, err
}

func createDB(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS snippets (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(100) NOT NULL,
			content TEXT NOT NULL,
			created DATETIME NOT NULL,
			expires DATETIME NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_snippets_created ON snippets(created);
	`

	_, err := db.Exec(query)

	return err
}
