// Bootstrap capability for the service

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	proto "test_service/protobuf/generated"
	"test_service/server"
)

// bootstrap the service during bootup
// initializes service config and logger module
// translates service config to its proto definition and return
func bootstrap(configPath string) (*proto.Config, error) {
	config, err := getConfig(configPath)
	if err != nil {
		log.Errorf("failed to read service config, path: %s, error: %s", configPath, err)
		return nil, err
	}

	// create an instance name for this service instance using a random number
	randomNum := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(256)
	instanceName := fmt.Sprintf("%s-%d", config.Service.Name, randomNum)

	protoConfig := &proto.Config{
		Service: &proto.ServiceConfig{
			Name:     config.Service.Name,
			FqdnOrIP: config.Service.FqdnOrIP,
			ApiPort:  config.Service.ApiPort,
			RpcPort:  config.Service.RpcPort,
		},
		Logging: &proto.LoggingConfig{
			LogDir:       config.Logging.LogDir,
			LogFile:      config.Logging.LogFile,
			LoggingLevel: config.Logging.LoggingLevel,
		},
		Host: &proto.HostConfig{
			Uuid:         uuid.New().String(),
			InstanceName: instanceName,
		},
	}

	return protoConfig, nil
}

// getConfig extracts config from service config file (yml) and provided flags
func getConfig(path string) (*server.Config, error) {
	var config server.Config
	if path == "" {
		return nil, fmt.Errorf("service config file path is empty")
	}

	log.Infof("reading yaml config from %q", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
