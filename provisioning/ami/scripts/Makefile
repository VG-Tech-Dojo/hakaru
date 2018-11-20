# provisioning (run only in building ami)
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

/root/hakaru:
	mkdir -p /root/hakaru

/root/hakaru/team_name.txt: deploy/team_name.txt /root/hakaru
	cp team_name.txt /root/hakaru/team_name.txt

/root/hakaru/Makefile: deploy/Makefile /root/hakaru/team_name.txt /root/hakaru
	cp deploy/Makefile /root/hakaru/Makefile

app: /root/hakaru/Makefile
