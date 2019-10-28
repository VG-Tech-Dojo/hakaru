export AWS_PROFILE        ?= sunrise201911
export AWS_DEFAULT_REGION ?= ap-northeast-1

KAKERU_ARTIFACTS := $(AWS_PROFILE)-kakeru-artifacts

ARTIFACTS_COMMIT ?= latest

.PHONY: upload clean

_deploy: scenarios/*.template kakeru.mk *.py keygen.sh
	-rm -rf _deloy
	-mkdir _deploy
	cp -r scenarios _deploy/
	cp keygen.sh _deploy/
	cp kakeru.mk _deploy/Makefile
	cp *.py _deploy/

kakeru.zip: _deploy
	cd _deploy && zip -r ../kakeru.zip ./*

upload: kakeru.zip
	aws s3 cp kakeru.zip s3://$(KAKERU_ARTIFACTS)/$(ARTIFACTS_COMMIT)/kakeru.zip

clean:
	-rm -rf kakeru.zip _deploy
