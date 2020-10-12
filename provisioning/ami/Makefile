# packer wrapper
#

export AWS_PROFILE        ?= sunrise2020-y
export AWS_DEFAULT_REGION := ap-northeast-1

#ARTIFACTS_COMMIT ?= $(shell git rev-parse HEAD)
ARTIFACTS_COMMIT := latest
ARTIFACTS_BUCKET := $(AWS_PROFILE)-hakaru-artifacts

# https://hub.docker.com/r/hashicorp/packer/tags/
PACKER_VERSION := 1.6.4

PACKER_IMAGE := sunrise2020/packer:$(PACKER_VERSION)

.PHONY: all clean build docker

TO := hakaru

all: __require_val clean docker build

__require_val:
ifeq (,$(shell ls $(TO).json 2>/dev/null))
	$(error "missing template file: $(TO).json")
endif

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
		$(PACKER_IMAGE) inspect $(TO).json
	docker run --rm -it \
		-v ~/.aws:/root/.aws \
		-e TZ=Asia/Tokyo \
		-e AWS_PROFILE=$(AWS_PROFILE) \
		-e AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
		-v $(CURDIR):/work \
		-w /work \
		$(PACKER_IMAGE) build --only=amazon-ebs -var ARTIFACTS_COMMIT=$(ARTIFACTS_COMMIT) $(TO).json

docker: Dockerfile
	docker build -t $(PACKER_IMAGE) --build-arg VERSION=$(PACKER_VERSION) -f $< .
