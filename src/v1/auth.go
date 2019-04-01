package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// Token desc
type Token struct {
	Token string `json:"token"`
}

// Auth desc
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *server) AuthRoutes() *chi.Mux {

	router := chi.NewRouter()
	router.Get("/verify", s.Verify)
	router.Post("/authorize", s.Authorize)

	return router
}

func (s *server) Authorize(w http.ResponseWriter, r *http.Request) {

	var a Auth

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, nil)
		return
	}

	err = json.Unmarshal(b, &a)

	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, nil)
		return
	}

	if validErrs := a.validate(); len(validErrs) > 0 {
		render.Status(r, 400)
		render.JSON(w, r, nil)
		return
	}

	//
	render.JSON(w, r, Token{Token: "token123446556"})

}

func (s *server) Verify(w http.ResponseWriter, r *http.Request) {

	authorization := r.Header.Get("Authorization")

	if authorization == "" {
		render.Status(r, 401)
		render.JSON(w, r, nil)
		return
	}

	a := Token{
		Token: authorization,
	}

	render.JSON(w, r, a)
}

func (a *Auth) validate() url.Values {

	errs := url.Values{}

	// check if the title empty
	if a.Username == "" {
		errs.Add("username", "The username field is required!")
	}

	if a.Password == "" {
		errs.Add("password", "The password field is required!")
	}

	return errs
}
