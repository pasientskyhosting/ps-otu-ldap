package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Routes explained
func (s *server) Routes() *chi.Mux {

	s.router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	s.router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/group", s.GroupRoutes())
		r.Mount("/api/user", s.UserRoutes())
		r.Mount("/api/ping", s.PingRoutes())
		r.Mount("/api/auth", s.AuthRoutes())
	})

	return s.router
}
