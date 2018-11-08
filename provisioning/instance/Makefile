all: clean awslogs deploy

clean:
	@echo 'do nothing'

deploy: /etc/systemd/system/hakaru.service
	$(MAKE) start

/opt/hakaru/bin:
	mkdir -p /opt/hakaru/bin

/opt/hakaru/bin/hakaru: ../../hakaru /opt/hakaru/bin
	cp ../../hakaru /opt/hakaru/bin/hakaru
	chmod +x /opt/hakaru/bin/hakaru

/etc/systemd/system/hakaru.service: systemd/hakaru.service /opt/hakaru/bin/hakaru
	cp systemd/hakaru.service /etc/systemd/system/hakaru.service
	systemctl list-unit-files --type=service | grep hakaru
	systemctl enable hakaru

/etc/awslogs/awscli.conf: awslogs/awscli.conf
	cp awslogs/awscli.conf /etc/awslogs/awscli.conf

/etc/awslogs/awslogs.conf: awslogs/awslogs.conf
	cp wslogs/awslogs.conf /etc/awslogs/awslogs.conf

/var/awslogs/state:
	mkdir -p /var/awslogs/state:

awslogs: /etc/awslogs/awscli.conf /etc/awslogs/awslogs.conf /var/awslogs/state
	systemctl start awslogsd
	systemctl enable awslogsd

start:
	systemctl start hakaru
