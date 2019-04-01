package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Group struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (s *server) GroupRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Get("/{groupID}", GetAGroup)
	router.Delete("/{groupID}", DeleteGroup)
	router.Post("/", CreateGroup)
	router.Get("/", GetAllGroups)

	return router
}

func GetAGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	groups := Group{
		Slug:  groupID,
		Title: "Hello world",
		Body:  "Heloo world from planet earth",
	}
	render.JSON(w, r, groups) // A chi router helper for serializing and returning json
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "Deleted TODO successfully"
	render.JSON(w, r, response) // Return some demo response
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "Created TODO successfully"
	render.JSON(w, r, response) // Return some demo response
}

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups := []Group{
		{
			Slug:  "slug",
			Title: "Hello world",
			Body:  "Heloo world from planet earth",
		},
	}
	render.JSON(w, r, groups) // A chi router helper for serializing and returning json
}
