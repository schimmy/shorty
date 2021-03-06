SHELL := /bin/bash
PKG := github.com/schimmy/shorty
PKGS := $(PKG) $(PKG)/db
EXECUTABLE := shorty

.PHONY: test $(PKGS) build clean

test: $(PKGS)

clean:
	rm -f $(GOPATH)/src/$(PKG)/build/$(EXECUTABLE)

build: clean
	go build -o build/$(EXECUTABLE) $(PKG)

$(PKGS):
ifeq ($(LINT),1)
	golint $(GOPATH)/src/$@*/**.go
endif
	go get -d -t $@
	go test $@ -test.v
