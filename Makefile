.PHONY : default deps test build image docs

export VERSION = 0.1.0-beta1

HARDWARE = $(shell uname -m)
OS := $(shell uname)

default: deps test build

deps:
	@echo "Configuring Last.Backend Dependencies"
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

test:
	@echo "Testing Last.Backend"
	@sh ./hack/run-coverage.sh

docs: docs/*
	@echo "Build Last.Backend Documentation"
	@sh ./hack/build-docs.sh

build:
	@echo "== Pre-building cli configuration"
	mkdir -p build/linux && mkdir -p build/darwin
	@echo "== Building Last.Backend CLI"
	@bash ./hack/build-cross.sh cli

install:
	@echo "== Install binaries"
	@bash ./hack/install-cross.sh

image:
	@echo "== Pre-building configuration"
	@sh ./hack/build-images.sh

swagger-spec:
	@echo "== Generating Swagger spec for Last.Backend API"
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -b ./cmd/kit -m -o ./swagger.json
