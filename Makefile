SHELL := /bin/bash
PKG := $(shell go list)
EXECUTABLE := $(basename $(PKG))
PKGS := $(shell go list ./... | grep -v /vendor)
READMES := $(foreach pkg,$(PKGS),$(pkg)/README.md)

GOVERSION := $(shell go version | grep 1.5)
ifeq "$(GOVERSION)" ""
  $(error must be running Go version 1.5)
endif
export GO15VENDOREXPERIMENT=1

.PHONY: test $(PKGS) clean doc vendor

all: build test

test: $(PKGS)

GOLINT := $(GOPATH)/bin/golint
$(GOLINT):
	@go get github.com/golang/lint/golint

GODEP := $(GOPATH)/bin/godep
$(GODEP):
	@go get -u github.com/tools/godep

GODOCDOWN := $(GOPATH)/bin/godocdown
$(GODOCDOWN):
	@go get github.com/robertkrimen/godocdown/godocdown

$(PKGS): $(GOLINT)
	@echo "FORMATTING..."
	@gofmt -w=true $(GOPATH)/src/$@/*.go
	@echo "LINTING..."
	@$(GOLINT) $(GOPATH)/src/$@/*.go
	@echo ""
ifeq ($(COVERAGE),1)
	@go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	@go tool cover -html=$(GOPATH)/src/$@/c.out
else
	@echo "TESTING..."
	@go test -v $@
	@echo ""
endif


build: $(PKGS)
	go build -o build/$(EXECUTABLE) $(PKG)

clean:
	rm -rf build

doc: $(READMES)

$(READMES):
	$(MAKE) $(GOPATH)/src/$@

%/README.md: PATH := $(PATH):$(GOPATH)/bin
%/README.md: %/*.go $(GODOCDOWN)
	$(GODOCDOWN) $(shell dirname $@) > $@

vendor: $(GODEP)
	$(GODEP) save $(PKGS)
	find vendor/ -path '*/vendor' -type d | xargs -IX rm -r X # remove any nested vendor directories
