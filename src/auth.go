package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
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

	err = s.lc.LDAPAuthentication(a)

	if err != nil {
		log.Printf("User %s denied LDAP login: %s \n", a.Username, err)
		render.Status(r, 401)
		render.JSON(w, r, nil)
		return
	}

	_, tokenString, _ := s.token.Encode(jwt.MapClaims{"user_id": a.Username, "exp": jwtauth.ExpireIn(60 * time.Minute)})

	render.JSON(w, r, Token{Token: tokenString})

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
