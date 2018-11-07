.PHONY: all install dep-ensure dep-ensure-update imports fmt test run build clean upload

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

clean:
	rm -rf hakaru

# deployment

artifacts.tgz: hakaru tools db provisioning/instance
	tar czf $@ $^

# FIXME: これを変えよう
YOUR_TEAMNAME := sunrise2018

export AWS_PROFILE        ?= sunrise2018
export AWS_DEFAULT_REGION := ap-northeast-1

ARTIFACTS_BUCKET := $(YOUR_TEAMNAME)-hakaru-artifacts

# ci からアップロードできなくなった場合のターゲット
upload: clean artifacts.tgz
	aws s3 cp artifacts.tgz s3://$(ARTIFACTS_BUCKET)/latest/artifacts.tgz
	aws s3 cp artifacts.tgz s3://$(ARTIFACTS_BUCKET/$$(git rev-parse HEAD)/artifacts.tgz
