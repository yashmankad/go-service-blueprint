// Bootstrap capability for the service

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	proto "test_service/protobuf/generated"
	"test_service/server"
)

// globals
var (
	logLevels = map[string]log.Level{
		"trace": log.TraceLevel,
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
		"fatal": log.FatalLevel,
		"panic": log.PanicLevel,
	}
)

// bootstraps the service during bootup
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

	//logger := getLogger(config)
	//return config, logger, nil

	if err := configureLogger(config); err != nil {
		log.Errorf("failed to initialize logger")
		return nil, err
	}

	protoConfig := &proto.Config{
		Service: &proto.ServiceConfig{
			Name:     config.Service.Name,
			FqdnOrIP: config.Service.FqdnOrIP,
			ApiPort:  config.Service.ApiPort,
			RpcPort:  config.Service.RpcPort,
		},
		Host: &proto.HostConfig{
			Uuid:         uuid.New().String(),
			InstanceName: instanceName,
		},
	}

	return protoConfig, nil
}

// configure logging parameters for the service
func configureLogger(config *server.Config) error {
	// form the log dir and filepath
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to fetch user home dir, error: %v", err)
		return err
	}

	// if log directory does not exist, create it
	logDir := filepath.Join(homeDir, "logs")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Fatalf("failed to create log directory, error: %v", err)
			return err
		}
	}

	logFile := filepath.Join(logDir, config.Service.LogFile)
	fh, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Errorf("failed to log to file, using default stderr, error: %v", err)
		return err
	}

	log.SetOutput(fh)
	log.SetLevel(logLevels[config.Service.LoggingLevel])
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// toggle for debugging
	log.SetReportCaller(false)

	return nil
}

// extract config from request config file (yaml) and provided flags
func getConfig(path string) (*server.Config, error) {
	var config server.Config
	if path == "" {
		return nil, fmt.Errorf("service config file path is empty")
	}

	log.Printf("reading yaml config from %q", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
