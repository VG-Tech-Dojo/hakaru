# application deployment
#
# /root/hakaru/Makefile on ec2 instance
#
# NOTE: make 実行中に Makefile が更新されるのを避けるため分割している
#

ARTIFACTS_BUCKET ?= sunrise201911-hakaru-artifacts
ARTIFACTS_COMMIT ?= latest

all: clean /root/hakaru/app
	$(MAKE) -C /root/hakaru/app/provisioning/instance

clean:
	rm -rf /tmp/artifacts.tgz /root/hakaru/app

/tmp/artifacts.tgz:
	aws s3 cp s3://$(ARTIFACTS_BUCKET)/$(ARTIFACTS_COMMIT)/artifacts.tgz /tmp/artifacts.tgz

/root/hakaru/app: /tmp/artifacts.tgz
	mkdir -p /root/hakaru/app
	tar xzvf /tmp/artifacts.tgz -C /root/hakaru/app
