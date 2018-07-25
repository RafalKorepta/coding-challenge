GIT_HASH = $(shell git rev-parse HEAD | tr -d "\n")
VERSION = $(shell git describe --tags --always --dirty --match=*.*.*)
GO_PKGS= \
    github.com/RafalKorepta/coding-challenge

all: backend
.PHONY: all

init:
	go get -d -u github.com/golang/dep
	go get -u github.com/hairyhenderson/gomplate
	go get -u github.com/tebeka/go2xunit
	go get -u github.com/axw/gocov/...
	go get -u github.com/AlekSi/gocov-xml
	go get -u github.com/onsi/ginkgo/ginkgo
.PHONY: init

backend: lint test-backend build-linux-backend
.PHONY: backend

lint:
	golangci-lint run
.PHONY: lint

test-backend:
	go vet $(GO_PKGS)
	echo "mode: set" > coverage-all.out
	$(foreach pkg,$(GO_PKGS),\
		go test -v -race -coverprofile=coverage.out $(pkg) | tee -a test-results.out || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out || exit 1;)
	go tool cover -func=coverage-all.out
.PHONY: test-backend

build-container-locally: build-linux-backend build-container
.PHONY: build-container-locally

build-container:
	docker build -t rafalkorepta/coding-challenge-backend:local-latest .
.PHONY: build-container-locally

build-linux-backend:
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(GIT_HASH)" -o dist/portal-backend main.go
.PHONY: build-linux-backend

deploy:
	docker build -f Dockerfile -t $(DOCKER_USERNAME)/coding-challenge-backend:$(VERSION) .
	docker push $(DOCKER_USERNAME)/coding-challenge-backend:$(VERSION)
	docker logout
.PHONY: deploy