package main

import (
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	db     *sql.DB
	router *chi.Mux
}

func main() {

	dbFile := os.Getenv("DB_FILE")      // SQLite DB path
	port := os.Getenv("PORT")           // Server listen port
	listen := os.Getenv("LISTEN")       // Listen
	atoken := os.Getenv("API_TOKEN")    // API token for sync services
	ekey := os.Getenv("ENCRYPTION_KEY") // Encryption key

	if atoken == "" {
		log.Fatalf("env API_TOKEN not set!")
	}

	if ekey == "" {
		log.Fatalf("env ENCRYPTION_KEY not set!")
	}

	if dbFile == "" {
		dbFile = "/data/otu-ldap/otu.db"
	}

	if port == "" {
		port = "8080"
	}

	if listen == "" {
		listen = "localhost"
	}

	// Check dbFile exists
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		log.Fatalf("Could find db: %q", err)
	}

	dbconn, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		log.Fatalf("Could not open db: %q", err)
	}

	defer dbconn.Close()

	// New server
	s := server{
		router: chi.NewRouter(),
		db:     dbconn,
	}

	log.Printf("Started REST API on %s:%s with db %s\n", listen, port, dbFile)
	log.Fatal(http.ListenAndServe(listen+":"+port, s.Routes()))

}

// Do some auth stuff here
func checkToken(r *http.Request) bool {
	return true
}

// Do some auth stuff here
func checkAPIKey(r *http.Request) bool {
	return true
}

func (s *server) isLDAPAuthorized(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if !checkToken(r) {
			render.Status(r, 401)
			render.JSON(w, r, nil)
			return
		}

		h(w, r)
	}
}

func (s *server) isAPIKeyAuthorized(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if !checkAPIKey(r) {
			render.Status(r, 401)
			render.JSON(w, r, nil)
			return
		}

		h(w, r)
	}
}
