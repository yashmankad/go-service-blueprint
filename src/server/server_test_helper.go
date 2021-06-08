// Helper class to create service helper objects, mainly for Go unit tests

package server

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	proto "test_service/protobuf/generated"
)

// server helper object is mainly used by Go unit tests
// it contains the server object and other mock objects to test a service out
type ServerHelper struct {
	// server object
	server *Server
}

func NewServerTestHelper() (*ServerHelper, error) {
	serviceName := "test-service"

	// create an instance name for this service instance using a random number
	randomNum := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(256)
	instanceName := fmt.Sprintf("%s-%d", serviceName, randomNum)

	config := &proto.Config{
		Service: &proto.ServiceConfig{
			Name:     serviceName,
			FqdnOrIP: "localhost",
			ApiPort:  "8000",
			RpcPort:  "8001",
		},
		Host: &proto.HostConfig{
			Uuid:         uuid.New().String(),
			InstanceName: instanceName,
		},
	}

	server, err := NewServer(config)
	if err != nil {
		log.Error("failed to create server object for test_service")
		return nil, err
	}

	// XXX: initialize other mock objects
	return &ServerHelper{server: server}, nil
}

func (sh *ServerHelper) CloseServerTestHelper() {
	sh.server.Close()

	// XXX: close other mock interfaces if created during server init
}
