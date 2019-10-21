# provisioning (run only in building ami)
#
# パッケージのインストール等はここでやる
#
.PHONY: all cloud-config apt awscli ssh ulimit

all: cloud-config apt awscli ssh ulimit

cloud-config: cloud.cfg.d/99_defaults.cfg
	cp cloud.cfg.d/99_defaults.cfg /etc/cloud/cloud.cfg.d/
	chmod 0644 /etc/cloud/cloud.cfg.d/99_defaults.cfg
	chown -R root:root /etc/cloud/cloud.cfg.d/99_defaults.cfg

apt:
	apt-get -y update
	apt-get -y upgrade
	apt-get -y install tsung python3 zip unzip

awscli: apt
	apt-get -y install python3-pip
	pip3 install awscli boto3

ssh:
	-mkdir ~/.ssh
	chmod 700 ~/.ssh
	ssh-keygen -f ~/.ssh/id_rsa
	cat ~/.ssh/id_rsa.pub > ~/.ssh/authorized_keys
	echo 'StrictHostKeyChecking no' >> ~/.ssh/config
	chmod 700 ~/.ssh/authorized_keys
	chown -R ubuntu:ubuntu ~/.ssh

ulimit:
	echo 'fs.file-max = 65000' >> /etc/sysctl.conf
	echo '*    soft nofile 65536' >> /etc/security/limits.conf
	echo '*    hard nofile 65536' >> /etc/security/limits.conf
	echo 'root soft nofile 65536' >> /etc/security/limits.conf
	echo 'root hard nofile 65536' >> /etc/security/limits.conf
	echo 'session    required     pam_limits.so' >> /etc/pam.d/common-session
