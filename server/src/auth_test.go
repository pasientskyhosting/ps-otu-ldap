package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAuthVerify(t *testing.T) {

	a := newAPITest(t, "GET", "/v1/api/auth/verify", nil)
	defer a.tearDown(t)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest")))

	response := executeRequest(a.server, a.req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestAuthVerifyShouldFailWhenInvalidToken(t *testing.T) {

	a := newAPITest(t, "GET", "/v1/api/auth/verify", nil)
	defer a.tearDown(t)

	response := executeRequest(a.server, a.req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}
