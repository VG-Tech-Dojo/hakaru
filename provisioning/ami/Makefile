SHELL := /bin/bash

AWS_PROFILE ?= sunrise2024-z

#ARTIFACTS_COMMIT ?= $(shell git rev-parse HEAD)
ARTIFACTS_COMMIT := latest
ARTIFACTS_BUCKET := $(AWS_PROFILE)-hakaru-artifacts

# https://hub.docker.com/r/hashicorp/packer/tags/
PACKER_VERSION := 1.11.2-1

PACKER_IMAGE := sunrise2024/packer:$(PACKER_VERSION)

.PHONY: all clean build scripts.tgz docker

TO ?= hakaru

all: __require_val clean docker build

__require_val: $(TO).pkr.hcl

scripts.tgz:
	tar cvzf scripts.tgz -C scripts .

clean:
	-rm -rf *.tgz

packer = docker run --rm -it \
  --env-file <(aws-vault exec $(AWS_PROFILE) -- env | grep "AWS_" | grep -v "AWS_VAULT") \
  -e TZ=Asia/Tokyo \
  -v $(CURDIR):/work \
  -v $(CURDIR)/packer_config:/root/.config/packer \
  -w /work \
  $(PACKER_IMAGE)

inspect: __require_val
	$(packer) validate -syntax-only $(TO).pkr.hcl
	$(packer) inspect .

init: inspect
	$(packer) init -upgrade .

build: init scripts.tgz
	$(packer) build -only '$(TO).*' -var ARTIFACTS_COMMIT=$(ARTIFACTS_COMMIT) .

docker: Dockerfile
	docker build -t $(PACKER_IMAGE) --build-arg VERSION=$(PACKER_VERSION) -f $< .
