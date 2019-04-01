package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// User def
type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
}

// UserRoutes def
func (s *server) UserRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Get("/{Username}", s.GetUser)
	router.Delete("/{Username}", s.DeleteUser)
	router.Post("/", s.CreateUser)
	router.Get("/", s.GetAllUsers)

	return router
}

// GetUser def
func (s *server) GetUser(w http.ResponseWriter, r *http.Request) {

	var username string
	var password string
	var groupID int
	var groupName string

	username = chi.URLParam(r, "Username")

	sqlStatement := "SELECT username, password, group_id, group_name FROM users LEFT JOIN groups ON users.group_id = groups.id WHERE username=$1;"

	row := s.db.QueryRow(sqlStatement, username)

	switch err := row.Scan(&username, &password, &groupID, &groupName); err {

	case sql.ErrNoRows:
		log.Printf("No rows were returned for username %v!", username)
		render.Status(r, 404)
		render.JSON(w, r, nil)
		break
	default:
		render.JSON(w, r, User{
			Username:  username,
			Password:  password,
			GroupID:   groupID,
			GroupName: groupName,
		})
		break
	}
}

func (s *server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "Deleted successfully"
	render.JSON(w, r, response) // Return some demo response
}

func (s *server) CreateUser(w http.ResponseWriter, r *http.Request) {

	render.JSON(w, r, User{
		Username: "1",
		Password: "3",
		GroupID:  2,
	})
}

func (s *server) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	username := "testing"

	users := []User{
		{
			Username: username,
		},
	}
	render.JSON(w, r, users) // A chi router helper for serializing and returning json
}
