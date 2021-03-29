package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestGroupsGetAllGroups(t *testing.T) {

	a := newAPITest(t, "GET", "/api/v1/groups", nil)
	defer a.tearDown(t)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest", true)))
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

	a := newAPITest(t, "DELETE", "/api/v1/groups/apitemptest", nil)
	defer a.tearDown(t)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest", true)))
	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

}

func TestGroupsDeleteGroupWhenUnAuthorized(t *testing.T) {

	a := newAPITest(t, "DELETE", "/api/v1/groups/apitemptest", nil)
	defer a.tearDown(t)

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestGroupsDeleteGroupWhenNotAdmin(t *testing.T) {

	a := newAPITest(t, "DELETE", "/api/v1/groups/apitemptest", nil)
	defer a.tearDown(t)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest", false)))

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusForbidden, response.Code)

}

func TestGroupsGetAllGroupsShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest(t, "GET", "/api/v1/groups", nil)
	defer a.tearDown(t)

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestGroupsCreateGroup(t *testing.T) {

	// Skip while in CI
	skipCI(t)

	a := newAPITest(t, "POST", "/api/v1/ldap-groups/voip/groups", []byte(`{"group_name":"voip-random","description":"VoIP random test","lease_time":8600,"custom_properties":{"key_1":"hello","key_2":"world"}}`))
	defer a.tearDown(t)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest", true)))

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusForbidden, response.Code)

	fmt.Printf("body %s", response.Body.String())

}

func TestGroupsCreateGroupShouldFailWhenNotAdmin(t *testing.T) {

	a := newAPITest(t, "POST", "/api/v1/ldap-groups/voip/groups", nil)
	defer a.tearDown(t)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest", false)))

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusForbidden, response.Code)

}

func TestGroupsCreateGroupShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest(t, "POST", "/api/v1/ldap-groups/voip/groups", nil)
	defer a.tearDown(t)

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestGroupsGetAllGroupUsers(t *testing.T) {

	a := newAPITest(t, "GET", "/api/v1/groups/apitemptest/users", nil)
	defer a.tearDown(t)

	a.req.Header.Set("X-API-KEY", a.server.env.apiKey)

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var users = []User{}

	err := json.Unmarshal([]byte(response.Body.String()), &users)

	// handle parse error
	if err != nil {
		t.Errorf("Error while parsing body %s", response.Body.String())
	}

	// Check if error in body
	if users[0].Username == "" || users[0].Password == "" || users[0].GroupName != "apitemptest" {
		t.Errorf("Error with body: %+v", users)
	}

}

func TestGroupsGetAllGroupUsersShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest(t, "GET", "/api/v1/groups/voip/users", nil)
	defer a.tearDown(t)

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestGroupsGetAllGroupsInLDAPScope(t *testing.T) {

	a := newAPITest(t, "GET", "/api/v1/ldap-groups/apitemptest/groups", nil)
	defer a.tearDown(t)

	a.req.Header.Set("X-API-KEY", a.server.env.apiKey)

	response := executeRequest(a.server, a.req)

	// Check repsonse
	if checkResponseCode(t, http.StatusOK, response.Code) {

		var groups = []Group{}

		err := json.Unmarshal([]byte(response.Body.String()), &groups)

		// handle parse error
		if err != nil {
			t.Errorf("Error while parsing body %s", response.Body.String())
		}

		// Check if error in body
		if groups[0].GroupName == "" || groups[0].LdapGroupName == "" {
			t.Errorf("Error with body: %+v", groups)
		}
	}

}

func TestGroupsGetAllGroupsInLDAPScopeShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest(t, "GET", "/api/v1/ldap-groups/apitemptest/groups", nil)
	defer a.tearDown(t)

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestGroupsUpdateGroup(t *testing.T) {

	a := newAPITest(t, "PATCH", "/api/v1/groups/apitemptest", []byte(`{"group_name":"apitemptest-updatedname","description":"apitemptest updated desc","lease_time":363,"custom_properties":[{"key":"updated1","value":"updated2"},{"key":"hello","value":"2"}]}`))
	defer a.tearDown(t)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "kj", true)))

	response := executeRequest(a.server, a.req)

	// Check repsonse
	if checkResponseCode(t, http.StatusOK, response.Code) {

		var group Group

		err := json.Unmarshal([]byte(response.Body.String()), &group)

		// handle parse error
		if err != nil {
			t.Errorf("Error while parsing body %s", response.Body.String())
		}

		// Check if error in body
		if group.GroupName != "apitemptest-updatedname" || group.LeaseTime != 363 || group.Description != "apitemptest updated desc" {
			t.Errorf("Error with body: %+v", group)
		}
	}

}

func TestGroupsUpdateGroupShouldFailWhenUnAuthorized(t *testing.T) {

	a := newAPITest(t, "PATCH", "/api/v1/groups/apitemptest", nil)
	defer a.tearDown(t)

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

func TestGroupsUpdateGroupShouldFailWhenNotAdmin(t *testing.T) {

	a := newAPITest(t, "PATCH", "/api/v1/groups/apitemptest", nil)
	defer a.tearDown(t)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apitest", false)))

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusForbidden, response.Code)

}
