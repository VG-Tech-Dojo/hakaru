# application setup
#
# /root/hakaru/app/provisioning/instance/Makefile on ec2 instance
#

.PHONY: all clean deploy start healthcheck

all: clean deploy start healthcheck

clean:
	@echo 'do nothing'

deploy: /etc/systemd/system/hakaru.service
	$(MAKE) start
	$(MAKE) healtcheck

/opt/hakaru/bin:
	mkdir -p /opt/hakaru/bin

/opt/hakaru/bin/hakaru: ../../hakaru /opt/hakaru/bin
	cp ../../hakaru /opt/hakaru/bin/hakaru
	chmod +x /opt/hakaru/bin/hakaru

/etc/sysconfig/hakaru: sysconfig/hakaru
	cp sysconfig/hakaru etc/sysconfig/hakaru

/etc/systemd/system/hakaru.service: systemd/hakaru.service /opt/hakaru/bin/hakaru /etc/sysconfig/hakaru
	cp systemd/hakaru.service /etc/systemd/system/hakaru.service
	systemctl list-unit-files --type=service | grep hakaru
	systemctl enable hakaru

start:
	systemctl start hakaru

healthcheck:
	curl -v 'http://127.0.0.1:8081/hakaru?name=deploy&value=1'
