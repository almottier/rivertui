.PHONY: all build test clean install lint fmt

# Binary name
BINARY_NAME=rivertui

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
GOLINT=golangci-lint
GOFMT=$(GOCMD) fmt

all: clean build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./main.go

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

install:
	$(GOGET) -v ./...

lint:
	$(GOLINT) run ./...

fmt:
	$(GOFMT) ./...

deps:
	go mod download
	go mod tidy

dev: build
	./$(BINARY_NAME)
