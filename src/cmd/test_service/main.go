// main module for the service

package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"test_service/server"
)

// globals
var (
	// configPath of service config
	configPath = flag.String("c", "config.yaml", "Path to service config")
)

// main routine for the service
// initializes service config and server object and starts the service
func main() {
	flag.Parse()

	protoConfig, err := bootstrap(*configPath)
	if err != nil {
		log.Error("failed to bootstrap test_service")
		os.Exit(1)
	}

	log.Info("service bootstrap successful")

	server, err := server.NewServer(protoConfig)
	if err != nil {
		log.Error("failed to create server object for test_service")
		os.Exit(1)
	}

	server.ContextLogger.Info("server object created successfully")

	if err = server.Run(); err != nil {
		log.Error("failed to initialize test_service")
		os.Exit(1)
	}
}
