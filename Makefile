VER ?= latest
GOFILES := $(wildcard cmd/api-server/*.go)

ifeq ($(OS),Windows_NT)
	ARCH := $(PROCESSOR_ARCHITECTURE)
	GOOS := windows
else
	ARCH := $(shell uname -m)
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS := linux
	else ifeq ($(UNAME_S),Darwin)
		GOOS := darwin
	endif
endif

ifeq (`echo $(ARCH) | tr A-Z a-z`, arm64)
	GOARCH=arm64
else
	GOARCH=amd64
endif

build:
	GOOS=${GOOS} GOARCH=$(GOARCH) go build \
	  -ldflags "-s -w -X main.BuildAt=`date +%FT%T%z`" \
	  -o build/simple-http-server $(GOFILES)

run:
	build/simple-http-server

clean:
	rm -rf build/*

docker:
	docker build -t simple-http-server:$(VER) .

.PHONY: build run clean docker