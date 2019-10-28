# provisioning target (run only in ami building)
#
# パッケージのインストール等はここでやる
#
.PHONY: all cloud-config yum app

all: cloud-config yum app

cloud-config: cloud.cfg.d/99_defaults.cfg
	cp cloud.cfg.d/99_defaults.cfg /etc/cloud/cloud.cfg.d/
	chmod 0644 /etc/cloud/cloud.cfg.d/99_defaults.cfg
	chown -R root:root /etc/cloud/cloud.cfg.d/99_defaults.cfg

yum:
	yum -y update
	yum -y install https://dev.mysql.com/get/mysql-community-release-el7-5.noarch.rpm
	yum -y --enablerepo=mysql56-community install mysql-community-client
	yum -y install java-1.8.0-openjdk
	yum -y install https://s3.amazonaws.com/amazoncloudwatch-agent/amazon_linux/amd64/latest/amazon-cloudwatch-agent.rpm
	yum -y install https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/linux_amd64/amazon-ssm-agent.rpm

/root/hakaru:
	mkdir -p /root/hakaru

/root/hakaru/Makefile: Makefile /root/hakaru
	cp Makefile /root/hakaru/Makefile

app: /root/hakaru/Makefile

# application deployment
#
# /root/hakaru/Makefile on ec2 instance
#

ARTIFACTS_BUCKET ?= sunrise201911-hakaru-artifacts
ARTIFACTS_COMMIT ?= latest

deploy: clean /root/hakaru/app
	$(MAKE) -C /root/hakaru/app/provisioning/instance

clean:
	rm -rf /tmp/artifacts.tgz /root/hakaru/app

/tmp/artifacts.tgz:
	aws s3 cp s3://$(ARTIFACTS_BUCKET)/$(ARTIFACTS_COMMIT)/artifacts.tgz /tmp/artifacts.tgz

/root/hakaru/app: /tmp/artifacts.tgz
	mkdir -p /root/hakaru/app
	tar xzvf /tmp/artifacts.tgz -C /root/hakaru/app
