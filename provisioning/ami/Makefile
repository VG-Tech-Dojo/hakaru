export AWS_PROFILE        ?= sunrise2018
export AWS_DEFAULT_REGION := ap-northeast-1

.PHONY: all clean build

all:
	@more Makefile

scripts.tgz: scripts
	tar cvzf $@ -C $< .

clean:
	-rm -rf *.tgz

build: packer.json clean scripts.tgz
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
		hashicorp/packer build --only=amazon-ebs $<
