# application deployment
#
# /root/hakaru/Makefile on ec2 instance
#
# NOTE: make 実行中に Makefile を更新されるのを避けるため分割している
#

TEAMNAME := $(shell head -n1 team_name.txt)

ARTIFACTS_COMMIT ?= latest

all: /tmp/artifacts.tgz /root/hakaru/app
	tar xzvf /tmp/artifacts.tgz -C /root/hakaru/app
	$(MAKE) -C /root/hakaru/app

/tmp/artifacts.tgz:
	aws s3 cp s3://$(TEAMNAME)-hakaru-artifacts/$(ARTIFACTS_COMMIT)/artifacts.tgz /tmp/artifacts.tgz

/root/hakaru/app:
	mkdir -p /root/hakaru/app
