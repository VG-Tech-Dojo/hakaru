# FIXME: これを変えよう
YOUR_TEAMNAME := flashers

export AWS_PROFILE        ?= sunrise2018
export AWS_DEFAULT_REGION := ap-northeast-1

ARTIFACTS_COMMIT ?= $(shell git rev-parse HEAD)
ARTIFACTS_BUCKET := $(YOUR_TEAMNAME)-hakaru-artifacts

# https://hub.docker.com/r/hashicorp/packer/tags/
PACKER_VERSION := 1.3.2

.PHONY: all clean check_artifacts build

all: clean build

scripts.tgz: scripts
	tar cvzf $@ -C $< .

clean:
	-rm -rf *.tgz

check_artifacts:
	@echo 'artifacts checking ...'
	@echo 'latest:'
	aws s3 ls s3://$(ARTIFACTS_BUCKET)/latest/artifacts.tgz
	@echo "using:"
	aws s3 ls s3://$(ARTIFACTS_BUCKET)/$(ARTIFACTS_COMMIT)/artifacts.tgz

build: packer.json scripts.tgz check_artifacts
	docker run --rm -it \
		-v ~/.aws:/root/.aws \
		-e TZ=Asia/Tokyo \
		-e AWS_PROFILE=$(AWS_PROFILE) \
		-e AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
		-v $(CURDIR):/work \
		-w /work \
		hashicorp/packer:$(PACKER_VERSION) inspect $<
	docker run --rm -it \
		-v ~/.aws:/root/.aws \
		-e TZ=Asia/Tokyo \
		-e AWS_PROFILE=$(AWS_PROFILE) \
		-e AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
		-v $(CURDIR):/work \
		-w /work \
		hashicorp/packer:$(PACKER_VERSION) build --only=amazon-ebs -var ARTIFACTS_COMMIT=$(ARTIFACTS_COMMIT) $<
