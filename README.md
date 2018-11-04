# Payments API

[![CircleCI](https://circleci.com/gh/carlosroman/payments-api.svg?style=svg)](https://circleci.com/gh/carlosroman/payments-api)[![Coverage Status](https://coveralls.io/repos/github/carlosroman/payments-api/badge.svg?branch=master)](https://coveralls.io/github/carlosroman/payments-api?branch=master)


## Setup

The project requires the following:
* Golang (1.11+)
* Docker
* Docker-compose

And a clone of the project

```
$ git clone https://github.com/carlosroman/payments-api.git
```

## Building

Open a terminal in the project directory and run:

```
$ make build
```

The executable can then be found in `checkout_dir/target` and should be called `server`.

If you prefer a Docker image you can run the following to build it:

```
$ make docker-build
```

### Notes

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules).
You might have to set the environment variable `GO111MODULE` to `on` if you have checked out the project into you GOPATH.

## Running tests

The easiest way to run the tests is by running:

```
$ make test
```

This uses [Ginkgo](https://onsi.github.io/ginkgo/) so you may need to install it by running:

```
go get github.com/onsi/ginkgo/ginkgo
```

## Running the application

The simplest way to run application is using Docker and Docker-compose.
This is done by running the following command from the project directory:

```
$ make start
```

Once the application has spun up you can go to [http://localhost:3000](http://localhost:3000) where the API documentation can be found.
The above command is actually running the Docker Compose file found [here](deployments/docker-compose.yml).
This spins up a [PostgreSQL](https://www.postgresql.org/), 
the [Swagger UI](https://swagger.io/tools/swagger-ui/) (which loads the swagger file which can be found [here](api/swagger.yaml)) and
the application server.

To stop the services running use:

```
$ make stop
```
