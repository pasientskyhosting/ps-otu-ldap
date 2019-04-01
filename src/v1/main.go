package main

import (
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	db     *sql.DB
	router *chi.Mux
}

// Examples
// docker run -e DB_FILE=/app/db.db -e PORT=8080 -e LISTEN=localhost -v /etc/otu:/app
func main() {

	dbFile := os.Getenv("DB_FILE")
	port := os.Getenv("PORT")
	listen := os.Getenv("LISTEN")

	if dbFile == "" {
		dbFile = "./otu.db"
	}

	if port == "" {
		port = "8080"
	}

	if listen == "" {
		listen = "localhost"
	}

	dbconn, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		log.Fatalf("Could not open db: %q", err)
	}

	defer dbconn.Close()

	s := server{
		router: chi.NewRouter(),
		db:     dbconn,
	}

	log.Printf("Started REST API on %s:%s with db %s\n", listen, port, dbFile)
	log.Fatal(http.ListenAndServe(listen+":"+port, s.Routes()))

}
