package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestLDAPGetAllLDAPGroups(t *testing.T) {

	a := newAPITest(t, "GET", "/api/v1/ldap-groups", nil)
	defer a.tearDown(t)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "kj", true)))
	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var ldapGroups []LDAPGroup

	err := json.Unmarshal([]byte(response.Body.String()), &ldapGroups)

	// handle parse error
	if err != nil {
		t.Errorf("Error while parsing body")
	}

	// Check if error in body
	if ldapGroups[0].LDAPGroupName == "" {
		t.Errorf("Error with body: %+v", ldapGroups)
	}

}

func TestLDAPGetAllLDAPGroupsShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest(t, "GET", "/api/v1/ldap-groups", nil)
	defer a.tearDown(t)

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestLDAPGetAllLDAPGroupsShouldFailWhenNotAdmin(t *testing.T) {

	a := newAPITest(t, "GET", "/api/v1/ldap-groups", nil)
	defer a.tearDown(t)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest", false)))

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusForbidden, response.Code)

}
