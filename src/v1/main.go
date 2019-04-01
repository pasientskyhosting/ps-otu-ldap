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

// Examples
// docker run -e DB_FILE=/app/db.db -e PORT=8080 -e LISTEN=localhost -v /etc/otu:/app
func main() {

	dbFile := os.Getenv("DB_FILE")
	port := os.Getenv("PORT")
	listen := os.Getenv("LISTEN")
	atoken := os.Getenv("API_TOKEN")
	ekey := os.Getenv("ENCRYPTION_KEY")

	if atoken == "" {
		log.Fatalf("env API_TOKEN not set!")
	}

	if ekey == "" {
		log.Fatalf("env ENCRYPTION_KEY not set!")
	}

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
