.PHONY: all clean test run fmt build provisioning/instance/etc/sazabi-ksk-keys.json

GOOS=
GOARCH=
GOSRC=$(GOPATH)/src

all: install run

install: dep-ensure

dep-ensure: Gopkg.toml
	which dep || go get github.com/golang/dep/...
	dep ensure

dep-ensure-update: Gopkg.toml
	which dep || go get github.com/golang/dep/...
	dep ensure --update

Gopkg.toml:
	which dep || go get github.com/golang/dep/...
	dep init

imports:
	goimports -w .

fmt:
	gofmt -w .

test:
	go test -v -tags=unit $$(go list ./...)

run: main.go
	go run main.go

build: test
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build

