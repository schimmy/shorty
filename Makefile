SHELL := /bin/bash
PKG := github.com/Clever/pretty-self-hosted-url-shortener
PKGS := $(PKG)
EXECUTABLE := shortener

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
