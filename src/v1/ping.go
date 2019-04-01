package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Ping struct {
	Version string `json:"version"`
	Message string `json:"message"`
}

func (s *server) PingRoutes() *chi.Mux {

	router := chi.NewRouter()
	router.Get("/", s.Ping)
	return router
}

func (s *server) Ping(w http.ResponseWriter, r *http.Request) {
	ping := Ping{
		Version: "v0.0.1",
		Message: "pong",
	}
	render.JSON(w, r, ping)
}
