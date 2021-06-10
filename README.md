# go-service-blueprint

This repository is a blueprint to create a micro-service in Golang.
It is a service skeleton that can bootstrap, configure and enable key service sub-systems like API and RPC servers, cross-language support (IDL), logging, unit test framework and the ability to build and deploy the service on Docker containers and Kubernetes.

The following sections provide more details on these sub-systems.

## Building the Service

We have defined a top level ```Makefile``` that defines various targets to build all binaries, including protocol buffer definitions. A simple ```make``` would build everything. Additional targets are defined to run unit tests, format code and build Docker artifacts. Please refer to the Makefile source code for additional details.

## Service Configuration

The service can be configured through ```yml``` file located at ```config/*.yml```. This config consists of service details like API and RPC port numbers, DNS name, etc along with details on how to connect to other vital services like the datastore, async queues, KV store and so on. The configuration is read when the service is bootstrapped and gets applied to the server object of the service as part of initialization.

## API and RPC Server

API and RPC servers are created asynchronously as part of server bringup and initialization. We leverage [**Gin**](https://github.com/gin-gonic/gin) for routing REST API requests and [**gRPC**](https://grpc.io/) to offer the capability to talk to the service using RPCs.

## Logging

The framework leverages the popular [**logrus**](https://github.com/sirupsen/logrus) Go package for logging all service logs, events and requests to the directory and file requested in the service configuration. 

## Cross Language Support

We lean on defining all essential request/response objects in the form of [**protocol buffer**](https://github.com/protocolbuffers/protobuf) definitions. This helps maintain backward compatibility across service versions and offers cross-language support where clients consuming this service need not be implemented in Golang. The gRPC framework also runs on top of protobufs ensuring consistent IDL usage.

## Unit Test Framework

The blueprint uses ```go test``` to develop and run unit tests for all packages.

## Deployment

Like most micro-services today, the blueprint is developed keeping in mind that it will run in [**Docker**](https://www.docker.com/) containers and will be eventually orchestrated via [**Kubernetes**](https://kubernetes.io/). ```make docker``` builds all the necessary artifacts including the Docker image for the service. Kubernetes deployment and service definitions are located at ```deployment/```.

