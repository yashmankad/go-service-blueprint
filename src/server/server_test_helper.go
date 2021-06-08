// Helper class to create service helper objects, mainly for Go unit tests

package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	proto "test_service/protobuf/generated"
	"test_service/util"
)

// server helper object is mainly used by Go unit tests
// it contains the server object and other mock objects to test a service out
type ServerHelper struct {
	// server object
	server *Server

	// keep track of any go routines created as part of the test
	wg *sync.WaitGroup
}

func NewServerTestHelper(testObj *util.Test) (*ServerHelper, error) {
	serviceName := "test_service"

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
		Logging: &proto.LoggingConfig{
			LogDir:       testObj.TestDir,
			LogFile:      "test_service.log",
			LoggingLevel: "info",
		},
		Host: &proto.HostConfig{
			Uuid:         uuid.New().String(),
			InstanceName: instanceName,
		},
	}

	server, err := NewServer(config)
	if err != nil {
		log.Errorf("failed to create server object for test_service: %v", err)
		return nil, err
	}

	// run server in a go routine since it is a blocking call
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = server.Run(); err != nil {
			log.Errorf("failed to initialize test_service: %v", err)
			return
		}
	}()

	// wait for server to come up
	server.WaitForServerBootup(5 * time.Second)

	// ensure API server is Up
	waitForAPIServer(5 * time.Second)

	// XXX: initialize or mock other objects
	return &ServerHelper{
		server: server,
		wg:     &wg}, nil
}

func (sh *ServerHelper) CloseServerTestHelper() {
	sh.server.Close()
	sh.wg.Wait()

	// XXX: close other mock interfaces if created during server init
}

func waitForAPIServer(timeout time.Duration) error {
	log.Info("checking api server availability")
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://127.0.0.1:8000/v1/ping")
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("api server is Up")
			return nil
		}

		resp.Body.Close()
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("api server failed to come up (timedout")
}
