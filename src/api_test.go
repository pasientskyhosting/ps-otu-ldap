package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

var testServer *server

type apiTest struct {
	req    *http.Request
	server *server
	token  string
}

func newAPITest(method string, url string, body []byte) *apiTest {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("Could not create API test request %v", err)
	}

	s := newServer(
		newEnv(os.Getenv("DB_FILE"),
			os.Getenv("LISTEN"),
			"API_TOKEN",
			"ENCRYPTION_KEY",
			"tet123",
			os.Getenv("LDAP_BASE"),
			"localhost",
			"636",
			os.Getenv("LDAP_BIND_DN"),
			os.Getenv("LDAP_BIND_PASSWORD"),
		),
	)

	a := apiTest{
		req:    req,
		server: s,
	}

	return &a

}

func TestPing(t *testing.T) {

	a := newAPITest("GET", "/v1/api/ping", nil)

	response := a.executeRequest(a.req)

	a.checkResponseCode(t, http.StatusOK, response.Code)

	// if body := response.Body.String(); body != "[]" {
	// 	t.Errorf("Expected an empty array. Got %s", body)
	// }
}

func (a *apiTest) getToken() string {

	// Create temp token for api test
	_, ts, _ := a.server.token.Encode(jwt.MapClaims{"user_id": "apiTest", "exp": jwtauth.ExpireIn(1 * time.Minute)})

	return ts
}

func TestUsersCreateUser(t *testing.T) {

	a := newAPITest("POST", "/v1/api/users", []byte(`{"group_name":"macandcheese"}`))

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.getToken()))

	response := a.executeRequest(a.req)

	a.checkResponseCode(t, http.StatusCreated, response.Code)

	var u User

	err := json.Unmarshal([]byte(response.Body.String()), &u)

	// handle parse error
	if err != nil {
		t.Errorf("Error while parsing body")
	}

	// Check if error in body
	if u.Username == "" || u.Password == "" || u.GroupName == "" {
		t.Errorf("Error with body: %+v", u)
	}

}

func TestUsersGetAllUsers(t *testing.T) {

	a := newAPITest("GET", "/v1/api/users", nil)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.getToken()))

	response := a.executeRequest(a.req)

	a.checkResponseCode(t, http.StatusOK, response.Code)

	var users []User

	err := json.Unmarshal([]byte(response.Body.String()), &users)

	// handle parse error
	if err != nil {
		t.Errorf("Error while parsing body")
	}

	// Check if error in body
	if users[0].Username == "" || users[0].Password == "" || users[0].GroupName == "" {
		t.Errorf("Error with body: %+v", users)
	}

}

func (a *apiTest) checkResponseCode(t *testing.T, expected, actual int) {

	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func (a *apiTest) executeRequest(req *http.Request) *httptest.ResponseRecorder {

	rr := httptest.NewRecorder()

	// Bind routes and server http
	a.server.routes().ServeHTTP(rr, req)

	return rr
}
