package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

// Routes explained
func (s *server) routes() *chi.Mux {

	s.router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	// JWT routes
	s.router.Group(func(r chi.Router) {

		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(s.token))

		// Handle valid / invalid tokens.
		r.Use(jwtauth.Authenticator)

		r.Get("/api/v1/auth/verify", s.Verify)
		r.Get("/api/v1/users", s.GetAllUsers)
		r.Post("/api/v1/groups", s.CreateGroup)
		r.Post("/api/v1/groups/{GroupName}/users", s.CreateUser)
		r.Delete("/api/v1/groups/{GroupName}", s.DeleteGroup)
		r.Get("/api/v1/groups", s.GetAllGroups)
	})

	// Public routes
	s.router.Group(func(r chi.Router) {
		r.Post("/api/v1/auth/authorize", s.Authorize)
		r.Get("/api/v1/ping", s.Ping)
	})

	// API-KEY routes
	s.router.Group(func(r chi.Router) {
		r.Get("/api/v1/groups/{GroupName}/users", s.isAPIKeyAuthorized(s.GetAllGroupUsers))
	})

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "html")

	FileServer(s.router, "/", http.Dir(filesDir))

	return s.router
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

}

// Do some auth stuff here
func checkToken(r *http.Request) bool {
	return true
}

// Do some auth stuff here
func checkAPIKey(r *http.Request, s *server) bool {

	if r.Header.Get("X-API-KEY") == "" {
		return false
	}

	if s.env.apiKey == r.Header.Get("X-API-KEY") {
		return true
	}

	return false
}

func (s *server) isAPIKeyAuthorized(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if !checkAPIKey(r, s) {
			render.Status(r, 401)
			render.JSON(w, r, nil)
			return
		}

		h(w, r)
	}
}
