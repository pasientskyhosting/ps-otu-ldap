package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestUsersCreateUser(t *testing.T) {

	var u User

	// create user
	a := newAPITest("POST", "/v1/api/groups/apitemptest/users", nil)
	defer a.tearDown()

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest")))

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusCreated, response.Code)

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

func TestUsersCreateUserShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest("POST", "/v1/api/groups/voip/users", nil)
	defer a.tearDown()

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestUsersGetAllUsers(t *testing.T) {

	var users []User

	// Get user
	a := newAPITest("GET", "/v1/api/users", nil)
	defer a.tearDown()

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest")))

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusOK, response.Code)

	err := json.Unmarshal([]byte(response.Body.String()), &users)

	// handle parse error
	if err != nil {
		t.Errorf("Error while parsing body")
	}

	// Check if error in body
	if users[0].Username == "" || users[0].Password == "" || users[0].GroupName != "apitemptest" {
		t.Errorf("Error with body: %+v", users)
	}

}

func TestUsersGetAllUsersShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest("GET", "/v1/api/users", nil)
	defer a.tearDown()

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}
