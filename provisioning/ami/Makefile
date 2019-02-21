# packer wrapper
#

TEAMNAME := $(shell head -n1 ../../team_name.txt)

export AWS_PROFILE        ?= $(TEAMNAME)
export AWS_DEFAULT_REGION := ap-northeast-1

#ARTIFACTS_COMMIT ?= $(shell git rev-parse HEAD)
ARTIFACTS_COMMIT := latest
ARTIFACTS_BUCKET := $(TEAMNAME)-hakaru-artifacts

# https://hub.docker.com/r/hashicorp/packer/tags/
PACKER_VERSION := 1.3.2

.PHONY: all clean check_artifacts build

all: clean build

scripts/deploy/team_name.txt: ../../team_name.txt
	cp ../../team_name.txt scripts/deploy/team_name.txt

scripts.tgz: scripts scripts/deploy/team_name.txt
	tar cvzf scripts.tgz -C scripts .

clean:
	-rm -rf *.tgz scripts/deploy/team_name.txt

check_artifacts:
	@echo 'artifacts checking ...'
	@echo "use:"
	aws s3 ls s3://$(ARTIFACTS_BUCKET)/$(ARTIFACTS_COMMIT)/artifacts.tgz

build: scripts.tgz check_artifacts
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
