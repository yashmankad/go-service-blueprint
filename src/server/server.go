// Server implementation for the service
// This package exposes capability to initialize a new server object and
// initialize its connections to the datastore, kv store and also standup
// REST and RPC servers

package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	proto "test_service/protobuf/generated"
	"test_service/repository"
	"test_service/router"
)

// globals
var (
	// logLevels helps map the user's requested logging level to the logger used
	logLevels = map[string]log.Level{
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
		"fatal": log.FatalLevel,
		"panic": log.PanicLevel,
	}
)

// Server object for the service
// contains handlers to api/rpc server, db object, server config, logger, etc
type Server struct {
	// service instance configuration
	Config *proto.Config

	// api server object
	ApiSrvr *http.Server

	// rpc server object (RPCs are implemented through gRPC)
	RpcSrvr *grpc.Server

	// file handle to server's logs
	LogFileHandle *os.File

	// logger object to log with additional service context
	ContextLogger *log.Entry

	// lock to ensure server object modifications are thread safe
	serverLock sync.Mutex

	// Waitgroup ensures the service stays running till all underlying threads exit
	wg sync.WaitGroup

	// flag to indicate if service is up/running
	running bool

	// XXX: needed due to some quirky grpc behavior (https://github.com/grpc/grpc-go/issues/3794)
	proto.UnimplementedTestServiceRPCServer

	// repository object (includes conn object to the db/repo)
	Repository *repository.Repository

	// XXX: add connection objects for kvstores, queues, etc
}

// NewServer initializes a new server object
// the object also contains a context logger to log with additional service context
func NewServer(config *proto.Config) (*Server, error) {
	// print the config before using it to initialize the server
	log.Infof("service config: %v", config)

	serverObj := &Server{Config: config}
	if err := serverObj.configureLogger(); err != nil {
		log.Errorf("failed to initialize logger")
		return nil, err
	}

	return serverObj, nil
}

// Run the server. as part of the process we initialize connections to
// datastore, KV stores and also spin up a REST/RPC server
func (s *Server) Run() error {
	s.serverLock.Lock()
	defer s.serverLock.Unlock()

	// initialize server instance
	if err := s.initialize(); err != nil {
		s.ContextLogger.Errorf("failed to initialize server instance: %v", err)
		return err
	}

	// start api server. this is done in a background thread since it is a blocking call
	s.wg.Add(2)
	go s.goRunAPIServer()
	go s.goRunRPCServer()

	s.running = true
	s.ContextLogger.Info("server started successfully")

	// release the lock so that background threads don't starve (in case they need the lock too)
	s.serverLock.Unlock()
	s.wg.Wait()
	s.serverLock.Lock()

	// if control comes here, it means both rest and rpc servers have terminated
	// mark the server instance as not running
	s.running = false

	return nil
}

// Close a server connection gracefully and stop the service instance
// best effort...currently does not return any error
func (s *Server) Close() error {
	s.running = false
	// close/stop the api server
	// the context is informs the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.ApiSrvr.Shutdown(ctx); err != nil {
		s.ContextLogger.Error("failed to shutdown api server, error:", err)
	}

	// stop the grpc server
	s.RpcSrvr.Stop()

	// close db conn
	// gorm supports connection pooling so you should only close this connection
	// if all consumers are done with it
	if s.Repository != nil {
		sqlDB, _ := s.Repository.DbConn.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return nil
}

// WaitForServerBootup is a utility that helps clients wait before issuing requests against the server
func (s *Server) WaitForServerBootup(timeout time.Duration) error {
	currTime := time.Now()
	endTime := currTime.Add(timeout)
	for time.Now().Before(endTime) {
		if s.getServerStatus() {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("server bootup timed out")
}

// initialize will setup connections from the server to external systems like datastore, KV store, queues, etc
func (s *Server) initialize() error {
	// comment out DB initialization since the template is not pointing to a valid DB (init will fail)
	// START BLOCK
	// initialize db connection
	/*
		if err := s.initializeDbConn(); err != nil {
			s.ContextLogger.Errorf("failed to initialize repository connection: %v", err)
			return err
		}
	*/
	// END BLOCK

	// XXX: initialize other connection objects like kvstore and queues
	return nil
}

func (s *Server) initializeDbConn() error {
	repoConn, err := repository.NewRepository(s.Config.Datastore)
	if err != nil {
		return err
	}

	s.ContextLogger.Infof("repository connection initialized successfully")
	s.Repository = repoConn
	return nil
}

// goRunAPIServer initializes the server's REST API server in the form of a Go routine
func (s *Server) goRunAPIServer() {
	defer s.wg.Done()

	r, err := router.NewRouter(s.LogFileHandle, s.Repository, s.ContextLogger)
	if err != nil {
		s.ContextLogger.Error("failed to initialize api router")
		s.wg.Done()
		return
	}

	apiSrvr := &http.Server{
		Addr:    ":" + s.Config.Service.ApiPort,
		Handler: r,
	}

	// store the api handle before starting the server
	// since starting the server is a blocking call
	s.ApiSrvr = apiSrvr

	// start the api server
	if err := apiSrvr.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			s.ContextLogger.Fatalf("error running api server on port %s, err: %s",
				s.Config.Service.ApiPort, err)
		}

		return
	}
}

// goRunRPCServer initializes the server's RPC server in the form of a Go routine
func (s *Server) goRunRPCServer() {
	s.wg.Done()
	listener, err := net.Listen("tcp", "localhost:"+s.Config.Service.RpcPort)
	if err != nil {
		s.ContextLogger.Fatalf("failed to listen on rpc port: %v", err)
	}

	grpcServer := grpc.NewServer()
	s.RpcSrvr = grpcServer
	proto.RegisterTestServiceRPCServer(grpcServer, s)
	if err := grpcServer.Serve(listener); err != nil {
		s.ContextLogger.Infof("error running rpc server on port %s, err: %s",
			s.Config.Service.RpcPort, err)
		return
	}
}

// configureLogger initializes the logging parameters for the server/service
func (s *Server) configureLogger() error {
	// if log directory does not exist, create it
	logDir := s.Config.Logging.LogDir
	if logDir == "" {
		// form the log dir and filepath
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to fetch user home dir, error: %v", err)
			return err
		}

		logDir = filepath.Join(homeDir, "logs")
	}

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Fatalf("failed to create log directory, error: %v", err)
			return err
		}
	}

	logFile := filepath.Join(logDir, s.Config.Logging.LogFile)
	fh, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Errorf("failed to log to file, using default stderr, error: %v", err)
		return err
	}

	s.LogFileHandle = fh

	log.SetOutput(fh)
	log.SetLevel(logLevels[s.Config.Logging.LoggingLevel])
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetFormatter(&log.JSONFormatter{})

	// create a logger object that logs with added context (instanceName, ip, etc)
	contextLogger := log.WithFields(log.Fields{
		"fqdn":     s.Config.Service.FqdnOrIP,
		"instance": s.Config.Host.InstanceName,
	})

	s.ContextLogger = contextLogger

	// toggle for debugging
	log.SetReportCaller(false)

	return nil
}

// getServerStatus safely fetches current server status (running vs. otherwise)
func (s *Server) getServerStatus() bool {
	s.serverLock.Lock()
	status := s.running
	s.serverLock.Unlock()
	return status
}
