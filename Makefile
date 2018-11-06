.PHONY: all install dep-ensure dep-ensure-update imports fmt test run build

GOOS   ?=
GOARCH ?=
GOSRC  := $(GOPATH)/src

all: install run

install: dep-ensure

dep-ensure: Gopkg.toml
	test -x $(GOPATH)/bin/dep || go get github.com/golang/dep/...
	$(GOPATH)/bin/dep ensure -v -vendor-only=true

dep-ensure-update: Gopkg.toml
	test -x $(GOPATH)/bin/dep || go get github.com/golang/dep/...
	$(GOPATH)/bin/dep ensure -v --update

Gopkg.toml:
	test -x $(GOPATH)/bin/dep || go get github.com/golang/dep/...
	$(GOPATH)/bin/dep init

imports:
	goimports -w .

fmt:
	gofmt -w .

test:
	go test -v -tags=unit $$(go list ./... | grep -v '/vendor/')

run: main.go
	go run main.go

build: hakaru

hakaru: test
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@

artifacts.tgz: hakaru tools/
	tar czf $@ $^
