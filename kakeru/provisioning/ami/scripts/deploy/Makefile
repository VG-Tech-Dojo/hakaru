# application deployment
#
# /root/hakaru/Makefile on ec2 instance
#
# NOTE: make 実行中に Makefile を更新されるのを避けるため分割している
#

ARTIFACTS_BUCKET ?= sunrise201911-kakeru-artifacts
ARTIFACTS_COMMIT ?= latest


all: clean /opt/kakeru/app
clean:
	-rm -rf /tmp/kakeru.zip /opt/kakeru/app

/tmp/kakeru.zip:
	aws s3 cp s3://$(ARTIFACTS_BUCKET)/$(ARTIFACTS_COMMIT)/kakeru.zip /tmp/kakeru.zip


/opt/kakeru:
	-mkdir -p /opt/kakeru
	cp Makefile /opt/kakeru/Makefile

/opt/kakeru/app: /opt/kakeru /tmp/kakeru.zip
	mkdir -p /opt/kakeru/app
	unzip -d /opt/kakeru/app /tmp/kakeru.zip
