package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAuthVerify(t *testing.T) {

	a := newAPITest("GET", "/v1/api/auth/verify", nil)

	a.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.server.getToken(1, "apiTest")))

	response := executeRequest(a.server, a.req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestAuthVerifyShouldFailWhenInvalidToken(t *testing.T) {

	a := newAPITest("GET", "/v1/api/auth/verify", nil)

	response := executeRequest(a.server, a.req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}
