package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

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
			"API_KEY",
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

func checkResponseCode(t *testing.T, expected, actual int) {

	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(s *server, req *http.Request) *httptest.ResponseRecorder {

	rr := httptest.NewRecorder()

	// Bind routes and server http
	s.routes().ServeHTTP(rr, req)

	return rr
}
