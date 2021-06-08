// Contains server unit testcases
package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"test_service/util"
	"test_service/v1api"
)

func TestServer(test *testing.T) {
	// initialize test object
	testObj, err := util.TestInit("test-service")
	if err != nil {
		test.Errorf("failed to initialize test object: %v", err)
		return
	}

	// defer testObj.TestCleanup(test)

	serverHelper, err := NewServerTestHelper(testObj)
	if err != nil {
		test.Errorf("failed to initialize server helper object: %v", err)
		return
	}

	resp, err := http.Get("http://127.0.0.1:8000/v1/ping")
	if err != nil {
		test.Errorf("failed to issue REST call to test api server: %v", err)
		return
	}
	defer resp.Body.Close()

	var pingResponse v1api.PingResponse
	json.NewDecoder(resp.Body).Decode(&pingResponse)
	if pingResponse.Message != "pong" {
		test.Errorf("invalid response to 'ping' API REST")
		return
	}

	serverHelper.CloseServerTestHelper()
}
