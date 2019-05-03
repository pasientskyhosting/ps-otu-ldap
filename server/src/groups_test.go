package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestGroupsGetAllGroups(t *testing.T) {

	a := newAPITest("GET", "/v1/api/groups", nil)
	defer a.tearDown()

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest")))
	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var groups []Group

	err := json.Unmarshal([]byte(response.Body.String()), &groups)

	// handle parse error
	if err != nil {
		t.Errorf("Error while parsing body")
	}

	// Check if error in body
	if groups[0].GroupName == "" || groups[0].LdapGroupName == "" || groups[0].CreateTime == 0 {
		t.Errorf("Error with body: %+v", groups)
	}

}

func TestGroupsDeleteGroup(t *testing.T) {

	a := newAPITest("DELETE", "/v1/api/groups/apitemptest", nil)
	defer a.tearDown()

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest")))
	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

}

func TestGroupsDeleteGroupWhenUnAuthorized(t *testing.T) {

	a := newAPITest("DELETE", "/v1/api/groups/apitemptest", nil)
	defer a.tearDown()

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestGroupsGetAllGroupsShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest("GET", "/v1/api/groups", nil)
	defer a.tearDown()

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestGroupsCreateGroup(t *testing.T) {

	a := newAPITest("POST", "/v1/api/groups", []byte(`{"group_name": "proxy-sql","ldap_group_name": "proxy-sql","lease_time": 3600}`))
	defer a.tearDown()

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest")))

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var g Group

	err := json.Unmarshal([]byte(response.Body.String()), &g)

	// handle parse error
	if err != nil {
		t.Errorf("Error while parsing body %s", response.Body.String())
	}

	// Check if error in body
	if g.GroupName == "" || g.LdapGroupName == "" {
		t.Errorf("Error with body: %s", response.Body.String())
	}

}

func TestGroupsCreateGroupShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest("POST", "/v1/api/groups", []byte(`{"group_name": "proxy-sql","ldap_group_name": "proxy-sql","lease_time": 3600}`))
	defer a.tearDown()

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestGroupsGetAllGroupUsers(t *testing.T) {

	a := newAPITest("GET", "/v1/api/groups/apitemptest/users", nil)
	defer a.tearDown()

	a.req.Header.Set("X-API-KEY", a.server.env.apiKey)

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// Body is encrypted so this has to be read
	cipherKey := []byte(a.server.env.ekey)
	decryptedPayload, err := decryptHash(cipherKey, response.Body.String())

	if err != nil {
		log.Fatalf("error: %s", err)
	}

	var users = []User{}

	err = json.Unmarshal([]byte(decryptedPayload), &users)

	// handle parse error
	if err != nil {
		t.Errorf("Error while parsing body %s", decryptedPayload)
	}

	// Check if error in body
	if users[0].Username == "" || users[0].Password == "" || users[0].GroupName != "apitemptest" {
		t.Errorf("Error with body: %+v", users)
	}

}

func TestGroupsGetAllGroupUsersShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest("GET", "/v1/api/groups/voip/users", nil)
	defer a.tearDown()

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}