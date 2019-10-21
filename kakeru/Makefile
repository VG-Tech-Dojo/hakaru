export AWS_PROFILE ?= $(TEAMNAME)
export AWS_DEFAULT_REGION ?= ap-northeast-1

TEAMNAME := $(shell head -n1 ../team_name.txt)
KAKERU_ARTIFACTS := $(TEAMNAME)-kakeru-artifacts

ARTIFACTS_COMMIT ?= latest

.PHONY: upload clean

_deploy: scenarios/*.template kakeru.mk *.py ../team_name.txt keygen.sh
	-rm -rf _deloy
	-mkdir _deploy
	cp -r scenarios _deploy/
	cp keygen.sh _deploy/
	cp kakeru.mk _deploy/Makefile
	cp *.py _deploy/
	cp ../team_name.txt _deploy/

kakeru.zip: _deploy
	cd _deploy && zip -r ../kakeru.zip ./*

upload: kakeru.zip
	aws s3 cp kakeru.zip s3://$(KAKERU_ARTIFACTS)/$(ARTIFACTS_COMMIT)/kakeru.zip

clean:
	-rm -rf kakeru.zip _deploy
