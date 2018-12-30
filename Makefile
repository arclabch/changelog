# Copyright © 2018 ARClab, Lionel Riem - https://arclab.ch/
#
# Use of this source code is governed by the MIT license that can be found in
# the LICENSE file.

# Binary output and package
OUT := changelog

# Go package
PKG := github.com/arclabch/changelog

# Stop if Windows detected
ifeq ($(OS),Windows_NT) 
    $(error Windows is not supported.)
endif

# Get version from Git
VERSION := $(shell git describe --always --tags --long --dirty 2>/dev/null || echo undefined)

# Get the OS
KERNEL := $(shell sh -c 'uname 2>/dev/null || echo Unknown')

# Get the machine platform
PLATFORM := $(shell sh -c 'uname -m 2>/dev/null || echo Unknown')

# Set the build timestamp
TIMESTAMP := $(shell sh -c 'date +"%Y-%m-%dT%k:%M%z"')

# If no prefix, set it
ifeq ($(PREFIX),)

    # FreeBSD lands in /usr/local, /usr otherwise
    ifeq ($(KERNEL),FreeBSD)
        PREFIX := /usr/local
    else
        PREFIX := /usr
    endif
endif

# For test/lint/vet/fmt
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

# Where to put the work we do
ifeq ($(WORKDIR),)
    WORKDIR := work
endif

# Default target is build
all: build

# Normal build
build:
	go build -i -v -o ${WORKDIR}/${OUT} -ldflags="-w -s \
    -X main.version=${VERSION} -X main.kernel=${KERNEL} -X main.machine=${PLATFORM} \
    -X main.built=${TIMESTAMP}" ${PKG}

# Run Go Test on the package(s)
test:
	@go test -short ${PKG_LIST}

# Vet the package(s)
vet:
	@go vet ${PKG_LIST}

# Lint the files
lint:
	@for file in ${GO_FILES} ; do \
		golint $$file ; \
	done

# Clean, vet, lint and build
dev-build: clean vet lint build

# Build and install
install: build
	install -m 755 ${WORKDIR}/${OUT} ${DESTDIR}${PREFIX}/bin/

# Uninstall
uninstall:
	@rm ${DESTDIR}${PREFIX}/bin/${OUT}

# Clean
clean:
	@rm -rf ${WORKDIR}

# Format the files
fmt:
	@for file in ${GO_FILES} ;  do \
		gofmt -w $$file ; \
	done

# Check the files format and print a diff
fmt-diff:
	@for file in ${GO_FILES} ;  do \
		gofmt -d $$file ; \
	done

# Clean, vet, lint, build and run
run: dev-build
	./${WORKDIR}/${OUT}

.PHONY: build dev-build test vet lint install uninstall clean fmt fmt-diff run