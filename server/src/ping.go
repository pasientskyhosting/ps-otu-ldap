package main

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

// go build -ldflags "-X name1=value1 -X name2=value2" -o path/to/output
var (
	version string
	date    string
)

// Ping struct
type Ping struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	Message string `json:"message"`
}

func (s *server) Ping(w http.ResponseWriter, r *http.Request) {

	if testingMode {
		version = "test"
		date = time.Now().Format(time.RFC3339)
	}

	ping := Ping{
		Version: version,
		Date:    date,
		Message: "pong",
	}

	render.JSON(w, r, ping)
}
