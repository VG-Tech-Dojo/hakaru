# provisioning target (run only in ami building)
#
# パッケージのインストール等はここでやる
#
.PHONY: all cloud-config yum app

export AWS_DEFAULT_REGION := ap-northeast-1
export AWS_DEFAULT_OUTPUT := text
export AWS_PAGER          :=

all: cloud-config yum app

cloud-config: ../cloud.cfg.d/99_defaults.cfg
	cp $< /etc/cloud/cloud.cfg.d/
	chmod 0644 /etc/cloud/cloud.cfg.d/99_defaults.cfg
	chown -R root:root /etc/cloud/cloud.cfg.d/99_defaults.cfg

yum:
	yum -y update
	yum -y install https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm
	rpm --import https://repo.mysql.com/RPM-GPG-KEY-mysql-2022
	yum-config-manager --disable mysql80-community
	yum-config-manager --enable mysql57-community
	yum -y install mysql-community-client
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

ARTIFACTS_BUCKET ?= sunrise2023-z-hakaru-artifacts
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
