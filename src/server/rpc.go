package server

import (
	"context"

	proto "test_service/protobuf/generated"
)

// Ping rpc request handler
func (s *Server) Ping(ctx context.Context, request *proto.PingRequest) (*proto.PingResponse, error) {
	response := &proto.PingResponse{
		Message: "pong",
	}

	return response, nil
}
