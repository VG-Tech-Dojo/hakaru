# FIXME: これを変えよう
YOUR_TEAMNAME := flashers

.PHONY: all cloud-config yum app clean

all: cloud-config yum app clean

cloud-config: cloud.cfg.d/99_defaults.cfg
	cp cloud.cfg.d/99_defaults.cfg /etc/cloud/cloud.cfg.d/
	chmod 0644 /etc/cloud/cloud.cfg.d/99_defaults.cfg
	chown -R root:root /etc/cloud/cloud.cfg.d/99_defaults.cfg

yum:
	yum -y update
	yum -y install https://dev.mysql.com/get/mysql-community-release-el7-5.noarch.rpm
	yum -y --enablerepo=mysql56-community install mysql-community-client
	yum -y install java-1.8.0-openjdk

/root/deploy:
	mkdir -p /root/deploy

ARTIFACTS_COMMIT ?= latest
/tmp/artifacts.tgz:
	aws s3 cp s3://$(YOUR_TEAMNAME)-hakaru-artifacts/$(ARTIFACTS_COMMIT)/artifacts.tgz /tmp/artifacts.tgz

app: ~/deploy /tmp/artifacts.tgz
	tar xzvf /tmp/artifacts.tgz -C /root/deploy
	$(MAKE) -C /root/deploy/provisioning/instance

clean:
	rm -rf /var/lib/yum && sudo yum clean all
	rm -rf /tmp/files /home/ec2-user/files
	rm -f /home/ec2-user/etc /home/ec2-user/.ssh/authorized_keys
	rm -f /etc/ssh/*_key /etc/ssh/*_key.pub
	rm -f /etc/udev/rules.d/70-persistent-net.rules
