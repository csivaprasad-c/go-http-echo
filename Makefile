# Metadata about this makefile and position
MKFILE_PATH := $(lastword $(MAKEFILE_LIST))
CURRENT_DIR := $(patsubst %/,%,$(dir $(realpath $(MKFILE_PATH))))

BUILD_TARGET_ARCH=${BUILD_ARCH}
BUILD_TARGET_OS=${BUILD_OS}
BIN_NAME=httpecho

GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

ifeq ($(BUILD_TARGET_ARCH),)
	BUILD_TARGET_ARCH=$(shell uname -m)
	ifeq ($(BUILD_TARGET_ARCH),'aarch64')
		BUILD_TARGET_ARCH='arm64'
	endif
endif

ifeq ($(BUILD_TARGET_OS),)
	BUILD_TARGET_OS=$(shell uname | tr -s '[:upper:]' '[:lower:]')
endif

build:
	env GOOS=${BUILD_TARGET_OS} CGO_ENABLED=0 GOARCH=${BUILD_TARGET_ARCH} go build -o ${BIN_NAME} \
		-ldflags="-X 'go-http-echo/version.Name=${BIN_NAME}' -X 'go-http-echo/version.GitCommit=${GIT_COMMIT}'"  .
.PHONY: build

clean:
	rm -rf ${BIN_NAME}
.PHONY: clean

image: clean build
	@docker buildx build --platform ${BUILD_TARGET_OS}/${BUILD_TARGET_ARCH} --tag csivaprasadc/${BIN_NAME}:latest --push .
