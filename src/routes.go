package main

import (
	"net/http"

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

		r.Get("/v1/api/auth/verify", s.Verify)
		r.Post("/v1/api/users", s.CreateUser)
		r.Get("/v1/api/users", s.GetAllUsers)
		r.Post("/v1/api/groups", s.CreateGroup)
		r.Get("/v1/api/groups", s.GetAllGroups)
	})

	// Public routes
	s.router.Group(func(r chi.Router) {
		r.Post("/v1/api/auth/authorize", s.Authorize)
		r.Get("/v1/api/ping", s.Ping)
	})

	// API-KEY routes
	s.router.Group(func(r chi.Router) {
		r.Get("/v1/api/groups/{GroupName}/users", s.isAPIKeyAuthorized(s.GetAllGroupUsers))
	})

	// _, claims, _ := jwtauth.FromContext(r.Context())
	// log.Printf("GetAllUsers accessed by: %v \n", claims["user_id"])

	return s.router
}

// Do some auth stuff here
func checkToken(r *http.Request) bool {
	return true
}

// Do some auth stuff here
func checkAPIKey(r *http.Request) bool {
	return true
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
