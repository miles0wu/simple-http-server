VER ?= latest
ARCH ?= amd64
GOFILES := $(wildcard cmd/api-server/*.go)

ifeq ($(OS),Windows_NT)
	GOOS := windows
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS := linux
	else ifeq ($(UNAME_S),Darwin)
		GOOS := darwin
	endif
endif

build:
	GOOS=${GOOS} GOARCH=$(ARCH) go build \
	  -ldflags "-s -w -X main.BuildAt=`date +%FT%T%z`" \
	  -o build/simple-http-server $(GOFILES)

run:
	build/simple-http-server

clean:
	rm -rf build/*

docker:
	docker build -t simple-http-server:$(VER) .

.PHONY: build run clean docker