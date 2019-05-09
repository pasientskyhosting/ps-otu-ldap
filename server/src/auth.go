package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
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

// Create temp token for e.g api test
func (s *server) getToken(expire time.Duration, username string, isAdmin bool) string {

	_, ts, _ := s.token.Encode(jwt.MapClaims{"user_id": username, "exp": jwtauth.ExpireIn(time.Minute * expire), "is_admin": isAdmin, "display_name": username})

	return ts
}

func (s *server) getLDAPUser(r *http.Request) (LDAPUser, error) {

	var tokenString = strings.Split(r.Header.Get("Authorization"), " ")

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString[1], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.env.jwtSecret), nil
	})

	if err != nil {
		return LDAPUser{}, err
	}

	isAdminStr := fmt.Sprintf("%t", claims["is_admin"])
	var isAdmin bool

	if isAdminStr == "true" {
		isAdmin = true
	} else {
		isAdmin = false
	}

	return LDAPUser{
		Username:    fmt.Sprintf("%s", claims["user_id"]),
		Admin:       isAdmin,
		DisplayName: fmt.Sprintf("%s", claims["display_name"]),
	}, nil

}

func (s *server) Authorize(w http.ResponseWriter, r *http.Request) {

	var a Auth

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		render.Status(r, 401)
		render.JSON(w, r, nil)
		return
	}

	err = json.Unmarshal(b, &a)

	if err != nil {
		render.Status(r, 401)
		render.JSON(w, r, nil)
		return
	}

	if validErrs := a.validate(); len(validErrs) > 0 {
		render.Status(r, 401)
		render.JSON(w, r, nil)
		return
	}

	lu, err := s.lc.LDAPAuthentication(a)

	if err != nil {
		log.Printf("User %s denied LDAP login: %s \n", a.Username, err)
		render.Status(r, 401)
		render.JSON(w, r, nil)
		return
	}

	_, tokenString, _ := s.token.Encode(jwt.MapClaims{"user_id": lu.Username, "exp": jwtauth.ExpireIn(3600 * time.Minute), "is_admin": lu.Admin, "display_name": lu.DisplayName})

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
