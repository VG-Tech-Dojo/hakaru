# packer wrapper
#

export AWS_PROFILE        ?= sunrise201911
export AWS_DEFAULT_REGION := ap-northeast-1

#ARTIFACTS_COMMIT ?= $(shell git rev-parse HEAD)
ARTIFACTS_COMMIT := latest
ARTIFACTS_BUCKET := $(AWS_PROFILE)-hakaru-artifacts

# https://hub.docker.com/r/hashicorp/packer/tags/
PACKER_VERSION := 1.3.2

.PHONY: all clean build

all: clean build

scripts.tgz: scripts
	tar cvzf scripts.tgz -C scripts .

clean:
	-rm -rf *.tgz

build: scripts.tgz
	docker run --rm -it \
		-v ~/.aws:/root/.aws \
		-e TZ=Asia/Tokyo \
		-e AWS_PROFILE=$(AWS_PROFILE) \
		-e AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
		-v $(CURDIR):/work \
		-w /work \
		hashicorp/packer:$(PACKER_VERSION) inspect packer.json
	docker run --rm -it \
		-v ~/.aws:/root/.aws \
		-e TZ=Asia/Tokyo \
		-e AWS_PROFILE=$(AWS_PROFILE) \
		-e AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
		-v $(CURDIR):/work \
		-w /work \
		hashicorp/packer:$(PACKER_VERSION) build --only=amazon-ebs -var ARTIFACTS_COMMIT=$(ARTIFACTS_COMMIT) packer.json
