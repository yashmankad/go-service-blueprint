// Server implementation for the service
// This package exposes capability to initialize a new server object and
// initialize its connections to the datastore, kv store and also standup
// REST and RPC servers

package server

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	proto "test_service/protobuf/generated"
	"test_service/router"
)

// Server object for the service
// contains handlers to api/rpc server, db object, server config, logger, etc
type Server struct {
	// api server object
	ApiSrvr *http.Server

	// rpc server object (RPCs are implemented through gRPC)
	RpcSrvr *grpc.Server

	// XXX: needed due to some quirky grpc behavior (https://github.com/grpc/grpc-go/issues/3794)
	proto.UnimplementedTestServiceRPCServer

	// logger object to log with additional service context
	ContextLogger *log.Entry

	// lock to ensure server object modifications are thread safe
	serverLock sync.Mutex

	// Waitgroup ensures the service stays running till all underlying threads exit
	wg sync.WaitGroup

	// flag to indicate if service is up/running
	running bool

	// service instance configuration
	config *proto.Config

	// XXX: add connection objects for databases, kvstores, queues, etc
}

// initializes a new server object
// the object also contains a context logger to log with additional service context
func NewServer(config *proto.Config) (*Server, error) {
	serverObj := &Server{config: config}

	// create a logger object that logs with added context (instanceName, ip, etc)
	contextLogger := log.WithFields(log.Fields{
		"fqdn":     config.Service.FqdnOrIP,
		"instance": config.Host.InstanceName,
	})

	serverObj.ContextLogger = contextLogger
	return serverObj, nil
}

// starts the server. as part of the process we initialize connections to
// datastore, KV stores and also spin up a REST/RPC server
func (s *Server) Run() error {
	s.serverLock.Lock()
	defer s.serverLock.Unlock()

	// initialize server instance
	if err := s.initialize(); err != nil {
		s.ContextLogger.Error("failed to initialize server instance: %v", err)
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

// gracefully close a server connection and stop the service instance
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

	return nil
}

// initializes connections to datastore, KV store, queues, etc
func (s *Server) initialize() error {
	// XXX: initialize other connection objects like datastore, kvstore and queues
	return nil
}

// initializes and starts a REST API server
func (s *Server) goRunAPIServer() {
	r, err := router.NewRouter()
	if err != nil {
		s.ContextLogger.Error("failed to initialize api router")
		s.wg.Done()
		return
	}

	apiSrvr := &http.Server{
		Addr:    ":" + s.config.Service.ApiPort,
		Handler: r,
	}

	// store the api handle before starting the server
	// since starting the server is a blocking call
	s.ApiSrvr = apiSrvr

	// start the api server
	if err := apiSrvr.ListenAndServe(); err != nil {
		s.ContextLogger.Fatal("error running api server: %s\n", err)
		s.wg.Done()
		return
	}

	s.wg.Done()
}

func (s *Server) goRunRPCServer() {
	listener, err := net.Listen("tcp", "localhost:"+s.config.Service.RpcPort)
	if err != nil {
		s.ContextLogger.Fatal("failed to listen on rpc port: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterTestServiceRPCServer(grpcServer, s)
	if err := grpcServer.Serve(listener); err != nil {
		s.ContextLogger.Fatal("error running rpc server on port %s, err: %s\n", s.config.Service.RpcPort, err)
		s.wg.Done()
		return
	}
}
