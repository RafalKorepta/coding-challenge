# coding-challenge

## Build
[![codecov](https://codecov.io/gh/RafalKorepta/coding-challenge/branch/master/graph/badge.svg)](https://codecov.io/gh/RafalKorepta/coding-challenge)
[![Build Status](https://travis-ci.org/RafalKorepta/coding-challenge.svg?branch=master)](https://travis-ci.org/RafalKorepta/coding-challenge)
[![Go Report Card](https://goreportcard.com/badge/github.com/RafalKorepta/coding-challenge)](https://goreportcard.com/report/github.com/RafalKorepta/coding-challenge)

## Docker image
[![](https://images.microbadger.com/badges/image/rafalkorepta/coding-challenge-backend:v0.0.1-3-gedbcbe4.svg)](https://microbadger.com/images/rafalkorepta/coding-challenge-backend:v0.0.1-3-gedbcbe4 "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/rafalkorepta/coding-challenge-backend:v0.0.1-3-gedbcbe4.svg)](https://microbadger.com/images/rafalkorepta/coding-challenge-backend:v0.0.1-3-gedbcbe4 "Get your own version badge on microbadger.com")

## Email Service
This repository contains mail microservice written in golang language. The main purpose of 
such microservice is to serve robust and scalable mailing service. This service accepts new Emails
from frontend portal. Then it low-balanced between sendgrid service or Amazon Simple Email Service.

In production this service should request the TLS certificates from e.g. Hashicorp Vault.
That certificates will be mounted to the kubernetes POD. The service in the initialization process
will start REST service using generated TLS certificates.

Another option to serve Email service securely can be done by using Envoy proxy container inside 
the Email Service kubernetes POD. The exposing port will come from Envoy. All the requests then 
will be forwarded to backend container. Envoy will be proxy that terminate TLS connection. 
The isolation provided by container runtime should not allow for any other network connection 
beside this Envoy sidecar.

The TLS will be out of scope for this coding challenge.

## User guide

### Getting started
To start service on kubernetes cluster run `helm install deploy` 

If you don't have kubernetes just visit http://apwdokowdkpaokdp.com

## Developer guide

### Prerequisite
For setup golang development environment visit https://golang.org/doc/code.html.
When you finish the setup then run `make init` to download tools that are used for testing and building purpose.
If you want to recompile protobuf contract please follow this instruction 
[https://github.com/grpc-ecosystem/grpc-gateway#installation](https://github.com/grpc-ecosystem/grpc-gateway#installation)

### Build
To build the backend application run `make backend`. It will create in `dist` folder binary named `portal` 
that can be run in your computer architecture.

### Generate Protobuf contract
In the `pkg/api/email/v1alpha1` folder you can find `email.proto` from which the Makefile build gRPC service, 
gRPC gateway and swagger definition. If you want to change contract, then run `make` inside 
`pkg/api/email/v1alpha1` path.

### Swagger-ui
The swagger-ui was copied from [https://github.com/philips/grpc-gateway-example](https://github.com/philips/grpc-gateway-example).
If you want to generate the ui again, then run `hack/build-ui.sh` from root folder of this repo. The swagger-ui
will be available on [https://email-backend-service.com/swagger-ui](https://email-backend-service.com/swagger-ui).

### Certificates
In the `pkg/certs/local_certs` folder you can find generated certificates that should be not used in production. 
If you would like to generated new certificates please go the the `certs` folder and run `make` command.
For demo purpose and testing the `local_certs` will be used to secure service.

### Configuration
All the configuration can be found in `docs` folder. It can be provided via config file 
(default `.portal-backend`) or via programs flags or via environment variables 
(the same name as flags but with prefix `SMACC` and all upper case).

### Testing
For writing unit test please use ginkgo as a framework for writing behavioral tests.

### Create containers
To create the backend containers locally run `make build-container-locally`

### Docs
If you want to generate command line markdown documentation go to `cmd` folder and run `go generate .`.

### CI
This project uses Travis-CI as a continuous integration backend. It will do the following tasks:
- run unit tests
- build backend application
- create docker images
- deploy backend in heroku
