package server

import (
	"context"

	proto "test_service/protobuf/generated"
)

func (s *Server) Ping(ctx context.Context, request *proto.PingRequest) (*proto.PingResponse, error) {
	response := &proto.PingResponse{
		Message: "pong",
	}

	return response, nil
}
