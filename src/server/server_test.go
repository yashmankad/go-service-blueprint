// Contains server unit testcases
package server

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"google.golang.org/grpc"

	proto "test_service/protobuf/generated"
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

	defer testObj.TestCleanup(test)

	serverHelper, err := NewServerTestHelper(testObj)
	if err != nil {
		test.Errorf("failed to initialize server helper object: %v", err)
		return
	}

	// test server's api capability
	resp, err := http.Get("http://127.0.0.1:8000/v1/ping")
	if err != nil {
		test.Errorf("failed to issue REST call to test api server: %v", err)
		return
	}
	defer resp.Body.Close()

	var apiResponse v1api.PingResponse
	json.NewDecoder(resp.Body).Decode(&apiResponse)
	if apiResponse.Message != "pong" {
		test.Errorf("invalid response to API request")
		return
	}

	// test server's grpc capability
	grpcConn, err := grpc.Dial("127.0.0.1:8001", grpc.WithInsecure())
	if err != nil {
		test.Errorf("failed to create connection object for grpc client: %v", err)
		return
	}
	defer grpcConn.Close()

	grpcClient := proto.NewTestServiceRPCClient(grpcConn)
	rpcResponse, err := grpcClient.Ping(context.Background(), &proto.PingRequest{})
	if err != nil {
		test.Errorf("failed to issue rpc call to server")
		return
	}

	if rpcResponse.Message != "pong" {
		test.Errorf("invalid response to RPC request")
		return
	}

	serverHelper.CloseServerTestHelper()
}
