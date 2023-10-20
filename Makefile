#!/usr/bin/make

-include include.make

SHELL := /bin/bash

APP_NAME ?= testapp

DOCKER_HOST 		?= unix:///var/run/docker.sock
DOCKER_PLATFORM		?= linux/amd64

DOCKERFILE ?= ./deployment/docker/${APP_NAME}/Dockerfile
DOCKERFILE_CONTEXT := .

IMAGE_REGISTRY     		?= docker.io/anvlad11
IMAGE_REGISTRY_URI 		?= ${IMAGE_REGISTRY}/${APP_NAME}

APPS := dawn-api
VERSION := $(shell cat ./VERSION)
LD_FLAGS := "-X main.Version=${VERSION} -X main.Commit=${CI_COMMIT_SHA}"

COMMIT_REF_NAME    := $(shell git rev-parse --abbrev-ref HEAD)
CI_COMMIT_REF_NAME ?= ${COMMIT_REF_NAME}

TAG := ${CI_COMMIT_REF_NAME}

$(APPS):
	@echo "Building $@..."
	@CGO_ENABLED=0 go build --ldflags ${LD_FLAGS} -o bin/$@ ./cmd/$@/

build: $(APPS)

build-docker: ## Build the image
	@echo "Building docker container ${IMAGE_REGISTRY_URI}:${TAG}"
	@DOCKER_HOST=${DOCKER_HOST} docker build \
		--no-cache \
		--platform ${DOCKER_PLATFORM} \
		--file ${DOCKERFILE} \
		-t ${IMAGE_REGISTRY_URI}:${TAG} ${DOCKERFILE_CONTEXT}

run-local:
	go run ./cmd/app/ --config-path ./config.yaml

run-compose:
	docker compose -f ./deployment/compose/docker-compose.yaml up -d

generate-oapi-models:
	@go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@d516da7
	@mkdir -p ./pkg/testapp
	@oapi-codegen --config oapi-codegen.yaml docs/api/openapi.yaml

generate: generate-oapi-models