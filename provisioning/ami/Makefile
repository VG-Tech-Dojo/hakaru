export AWS_PROFILE        ?= sunrise2018
export AWS_DEFAULT_REGION := ap-northeast-1

.PHONY: all clean check_artifacts build

all: clean build

scripts.tgz: scripts
	tar cvzf $@ -C $< .

clean:
	-rm -rf *.tgz

GIT_COMMIT       := $(shell git rev-parse --short HEAD)
ARTIFACTS_COMMIT := $(shell aws s3 ls s3://sunrise2018-hakaru-artifacts/$(GIT_COMMIT)/artifacts.tgz >/dev/null 2>&1 && echo $(GIT_COMMIT) || echo latest)

check_artifacts:
	@echo 'artifacts checking ...'
	@echo 'latest:'
	@aws s3 ls s3://sunrise2018-hakaru-artifacts/latest/artifacts.tgz
	@echo "$(GIT_COMMIT):"
	@aws s3 ls s3://sunrise2018-hakaru-artifacts/$(GIT_COMMIT)/artifacts.tgz || echo 'missing.\n\nuse latest artifacts...'

build: packer.json scripts.tgz check_artifacts
	docker run --rm -it \
		-v ~/.aws:/root/.aws \
		-e TZ=Asia/Tokyo \
		-e AWS_PROFILE=$(AWS_PROFILE) \
		-e AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
		-v $(CURDIR):/work \
		-w /work \
		hashicorp/packer inspect $<
	docker run --rm -it \
		-v ~/.aws:/root/.aws \
		-e TZ=Asia/Tokyo \
		-e AWS_PROFILE=$(AWS_PROFILE) \
		-e AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
		-v $(CURDIR):/work \
		-w /work \
		hashicorp/packer build --only=amazon-ebs -var ARTIFACTS_COMMIT=$(ARTIFACTS_COMMIT) $<
