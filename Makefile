TEAMNAME := $(shell head -n1 team_name.txt)

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

build: test
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o hakaru

clean:
	rm -rf hakaru *.tgz

# deployment

ARTIFACTS_TMPDIR := $(shell mktemp -u)/artifacts

$(ARTIFACTS_TMPDIR): provisioning/instance
	$(MAKE) build GOOS=linux GOARCH=amd64
	mkdir -p $(ARTIFACTS_TMPDIR)
	cp -r provisioning/instance/* $(ARTIFACTS_TMPDIR)
	cp hakaru $(ARTIFACTS_TMPDIR)

artifacts.tgz: $(ARTIFACTS_TMPDIR)
	tar czf artifacts.tgz -C $(ARTIFACTS_TMPDIR) .
	rm -rf $(ARTIFACTS_TMPDIR)

export AWS_PROFILE        ?= sunrise2018
export AWS_DEFAULT_REGION := ap-northeast-1

ARTIFACTS_BUCKET := $(TEAMNAME)-hakaru-artifacts

# ci からアップロードできなくなった場合のターゲット
upload: clean artifacts.tgz
	aws s3 cp artifacts.tgz s3://$(ARTIFACTS_BUCKET)/latest/artifacts.tgz
	aws s3 cp artifacts.tgz s3://$(ARTIFACTS_BUCKET)/$$(git rev-parse HEAD)/artifacts.tgz
