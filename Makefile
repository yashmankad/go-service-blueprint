# initializing Go env variables
export GOPATH=$(CURDIR)
export GOBIN=$(CURDIR)/bin

export PATH=/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:$(GOBIN)

.DEFAULT_GOAL := all
.PHONY: all build fmt clean

UNAME=$(shell uname)

# collection of Go packages and binaries
GO_PACKAGE_LIST=$(shell cd $(CURDIR)/src; go list ./...)
GO_CMD_LIST=$(shell cd $(CURDIR)/src/cmd; go list ./...)
PROTOBUF_FILES=$(shell ls src/protobuf/*.proto)

# build Go binaries and saves them in $GOBIN
build:
	@cd $(CURDIR)/src; go mod tidy; cd $(CURDIR)
	@for cmd in $(GO_CMD_LIST); do \
		bin_path=$(GOBIN)/`/usr/bin/basename $$cmd`; \
		cd $(CURDIR)/src;CGO_ENABLED=0 GOOS=linux go build -o $$bin_path $$cmd; cd $(CURDIR); \
	done

build-osx:
	@cd $(CURDIR)/src; go mod tidy; cd $(CURDIR)
	@for cmd in $(GO_CMD_LIST); do \
		bin_path=$(GOBIN)/osx/`/usr/bin/basename $$cmd`; \
		cd $(CURDIR)/src;CGO_ENABLED=0 go build -o $$bin_path $$cmd; cd $(CURDIR); \
	done

# formats all Go files and packages using 'go fmt'
fmt:
	@for package in $(GO_PACKAGE_LIST); do \
		cd $(CURDIR)/src; go fmt $$package; cd $(CURDIR); \
	done

# builds protocol buffer definitions and generates Go bindings for them
protobuf:
	@go get -u github.com/golang/protobuf/protoc-gen-go
	@go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@mkdir -p src/protobuf/generated
	@for protofile in $(PROTOBUF_FILES); do \
		protoc -I=./src/protobuf --go_out=./src/protobuf/generated --go-grpc_out=./src/protobuf/generated $$protofile; \
	done

# builds everything - protocol buffers, Go binaries and formats the code
all:
	@$(MAKE) -w protobuf
	@$(MAKE) -w build
	@$(MAKE) -w build-osx
	@$(MAKE) -w fmt

# cleans up existing Go binaries if any
clean:
	-rm -rf $(GOBIN)
	-rm -rf $(GOBIN)/osx

# unit test target
test:
	@mkdir -p $(CURDIR)/testout
	@for package in $(GO_PACKAGE_LIST); do \
		cd $(CURDIR)/src; go test -v $$package -test_dir='$(CURDIR)/testout'; cd $(CURDIR); \
	done

docker:
	@$(MAKE) -w all
	@sudo docker build -t test_service .

	