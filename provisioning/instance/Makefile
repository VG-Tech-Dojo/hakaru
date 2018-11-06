all: clean deploy

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

start:
	systemctl start hakaru
