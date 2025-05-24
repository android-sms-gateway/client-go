# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_PATH=build
BINARY_NAME=$(BINARY_PATH)/smsgateway
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
clean: 
	$(GOCLEAN)
	rm -rf $(BINARY_PATH)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v .
	./$(BINARY_NAME)
deps:
	$(GOGET) -v ./...

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Linting
lint:
	golangci-lint run -v
