# coding-challenge

[![codecov](https://codecov.io/gh/RafalKorepta/coding-challenge/branch/master/graph/badge.svg)](https://codecov.io/gh/RafalKorepta/coding-challenge)
[![Build Status](https://travis-ci.org/RafalKorepta/coding-challenge.svg?branch=master)](https://travis-ci.org/RafalKorepta/coding-challenge)
[![Go Report Card](https://goreportcard.com/badge/github.com/RafalKorepta/coding-challenge)](https://goreportcard.com/report/github.com/RafalKorepta/coding-challenge)

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

### Build
To build the backend application run `make backend`. It will create in `dist` folder binary named `portal` 
that can be run in your computer architecture.

### Testing
For writing unit test please use ginkgo as a framework for writing behavioral tests.

### Create containers
To create the backend containers locally run `make build-container-locally`

### CI
This project uses Travis-CI as a continuous integration backend. It will do the following tasks:
- run unit tests
- build backend application
- create docker images
- deploy backend in heroku