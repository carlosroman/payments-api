.DEFAULT_GOAL := test

.PHONY: build test clean docker-build start stop restart dc-build fmt info

NS ?= carlosroman/payments-api
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
        -skipPackage vendor \
        --randomizeAllSpecs \
        --randomizeSuites \
        --failOnPending \
        --cover \
        --coverprofile=payments.coverprofile \
        --outputdir=target \
        --trace \
        --race \
        --compilers=$(GINKGO_COMPILERS) \
        ./...

clean:
	@rm -rf target

build: export CGO_ENABLED=0
build:
	@mkdir -p target
	@go build \
	    -a \
	    -installsuffix cgo \
	    -o ./target/server \
	    ./cmd/server/server.go

docker-build:
	@$(DOCKER) build \
	     -f ./build/docker/Dockerfile \
	     -t $(NS)/server:$(VERSION) .

start:
	@$(DOCKER_COMPOSE_CMD) up

stop:
	@$(DOCKER_COMPOSE_CMD) down

restart: stop start;

fmt:
	@go fmt ./...

dc-build:
	@$(DOCKER_COMPOSE_CMD) build

info:
	@go env
