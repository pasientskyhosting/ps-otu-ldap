package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {

	var ping Ping

	a := newAPITest("GET", "/v1/api/ping", nil)

	response := executeRequest(a.server, a.req)
	checkResponseCode(t, http.StatusOK, response.Code)

	err := json.Unmarshal([]byte(response.Body.String()), &ping)

	// handle parse error
	if err != nil {
		t.Errorf("Error while parsing body")
	}

	// Check if error in body
	if ping.Version == "" || ping.Message == "" {
		t.Errorf("Error with body: %+v", ping)
	}

}
