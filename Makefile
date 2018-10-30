.DEFAULT_GOAL := test

.PHONY: build test clean docker-build start stop

NS ?= payments-api
VERSION ?= latest

DOCKER ?= docker
DOCKER_COMPOSE_FILE := ./deployments/docker-compose.yml
DOCKER_COMPOSE ?= docker-compose
DOCKER_COMPOSE_CMD := $(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE)
GINKGO_COMPILERS ?= 2

test-clean: clean
	@mkdir -p target

test: test-clean
	@ginkgo \
        -r \
        --randomizeAllSpecs \
        --randomizeSuites \
        --failOnPending \
        --cover \
        --coverprofile=$(NS).coverprofile \
        --outputdir=target \
        --trace \
        --race \
        --compilers=$(GINKGO_COMPILERS) \
        ./...

clean:
	rm -rf target

build: export CGO_ENABLED=0
build:
	@mkdir -p target
	@go build \
	    -a \
	    -installsuffix cgo \
	    -o ./target/$(NS) \
	    ./cmd/server/server.go

docker-build:
	@$(DOCKER) build \
	     -f ./build/docker/Dockerfile \
	     -t $(NS)/server:$(VERSION) .

start:
	@$(DOCKER_COMPOSE_CMD) up

stop:
	@$(DOCKER_COMPOSE_CMD) down
