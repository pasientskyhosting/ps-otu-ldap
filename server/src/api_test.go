package main

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

type apiTest struct {
	req      *http.Request
	server   *server
	token    string
	tearDown func(t *testing.T)
}

func (s *server) deleteTestData() {

	// delete all previos apitest users and groups
	deleteGroups, err := s.db.Prepare("DELETE FROM groups WHERE create_by='apiTest'")

	if err != nil {
		log.Fatalf("ERROR: preparing delete statement %+v", err)
	}

	deleteUsers, err := s.db.Prepare("DELETE FROM users WHERE create_by='apiTest'")

	if err != nil {
		log.Fatalf("ERROR: preparing delete statement %+v", err)
	}

	_, err = deleteGroups.Exec()

	if err != nil {
		// handle this error better than this
		log.Fatalf("ERROR: Deleting apiTest groups from DB  %+v", err)
	}

	_, err = deleteUsers.Exec()

	if err != nil {
		// handle this error better than this
		log.Fatalf("ERROR: Deleting apiTest users from DB  %+v", err)
	}

}

func (s *server) insertTestData() {

	// first delete anything
	s.deleteTestData()

	insertGroup, err := s.db.Prepare("INSERT INTO groups (id, group_name, ldap_group_name, lease_time, deleted, create_by, create_time) values('-1','apitemptest','apitemptest',720,0,'apiTest','1')")

	if err != nil {
		log.Fatalf("ERROR: preparing delete statement %+v", err)
	}

	_, err = insertGroup.Exec()

	if err != nil {
		// handle this error better than this
		log.Fatalf("ERROR: Insert apiTest test groups%+v", err)
	}

	insertUser, err := s.db.Prepare("INSERT INTO users (username, password, group_id, expire_time, create_by, create_time) values('apiTest-apitemptest-testblabla','Tj9cfKOMV2mvCB2-ozftPHKq6SjrdJjbXbcscA==',-1,$1,'apiTest','1');")

	if err != nil {
		log.Fatalf("ERROR: preparing delete statement %+v", err)
	}

	expires := time.Now().Unix() + 3600

	_, err = insertUser.Exec(expires)

	if err != nil {
		// handle this error better than this
		log.Fatalf("ERROR: Insert apiTest test user%+v", err)
	}

}

func (s *server) setupTest(t *testing.T) func(t *testing.T) {

	v := reflect.ValueOf(t)
	name := v.FieldByName("name")

	log.Printf("Setup test %s", name)
	s.insertTestData()

	return func(t *testing.T) {

		log.Printf("Teardown test %s", name)
		s.deleteTestData()

	}

}

func newAPITest(t *testing.T, method string, url string, body []byte) *apiTest {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		log.Fatalf("Could not create API test request %v", err)
	}

	s := newServer(
		newEnv(os.Getenv("DB_FILE"),
			os.Getenv("LISTEN"),
			os.Getenv("API_KEY"),
			os.Getenv("ENCRYPTION_KEY"),
			os.Getenv("JWT_SECRET"),
			os.Getenv("LDAP_BASE"),
			os.Getenv("LDAP_SERVER"),
			os.Getenv("LDAP_PORT"),
			os.Getenv("LDAP_BIND_DN"),
			os.Getenv("LDAP_BIND_PASSWORD"),
		),
	)

	tearDown := s.setupTest(t)

	a := apiTest{
		req:      req,
		server:   s,
		tearDown: tearDown,
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

	// Start rand seed
	rand.Seed(time.Now().UnixNano())

	// Bind routes and server http
	s.routes().ServeHTTP(rr, req)

	return rr
}
