# instance setup
#
# /root/hakaru/app/Makefile on ec2 instance
#

# artifacts.tgz が /root/hakaru/app 展開されて、以下のようになる
#.
#├── Makefile
#├── amazon-cloudwatch-agent
#│   └── amazon-cloudwatch-agent.json
#├── hakaru
#├── sysconfig
#│   └── hakaru
#└── systemd
#    └── hakaru.service

.PHONY: all clean hakaru-daemon amazon-cloudwatch-agent healthcheck

all: clean amazon-cloudwatch-agent hakaru-daemon healthcheck

clean:
	-systemctl stop hakaru
	-systemctl stop amazon-cloudwatch-agent
	-rm -f /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.toml

/opt/hakaru/bin:
	mkdir -p /opt/hakaru/bin

/opt/hakaru/bin/hakaru: hakaru /opt/hakaru/bin
	cp ../../hakaru /opt/hakaru/bin/hakaru
	chmod +x /opt/hakaru/bin/hakaru

/etc/sysconfig/hakaru: sysconfig/hakaru
	cp sysconfig/hakaru /etc/sysconfig/hakaru

/etc/systemd/system/hakaru.service: systemd/hakaru.service /opt/hakaru/bin/hakaru /etc/sysconfig/hakaru
	cp systemd/hakaru.service /etc/systemd/system/hakaru.service
	systemctl daemon-reload
	systemctl list-unit-files --type=service | grep hakaru

hakaru-daemon: /etc/systemd/system/hakaru.service
	systemctl start hakaru
	systemctl enable hakaru

healthcheck:
	curl -v 'http://127.0.0.1:8081/hakaru?name=deploy&value=1'

# https://docs.aws.amazon.com/ja_jp/AmazonCloudWatch/latest/monitoring/CloudWatch-Agent-Configuration-File-Details.html
/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json: amazon-cloudwatch-agent/amazon-cloudwatch-agent.json
	cp amazon-cloudwatch-agent/amazon-cloudwatch-agent.json /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json

/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.toml: /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json
	/opt/aws/amazon-cloudwatch-agent/bin/config-translator \
	    --mode ec2 \
	    --os linux \
	    --input /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json \
	    --output /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.toml \
	    --config /opt/aws/amazon-cloudwatch-agent/etc/common-config.toml
	/opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent \
	    -schematest \
	    -config /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.toml

amazon-cloudwatch-agent: /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.toml
	systemctl restart amazon-cloudwatch-agent
	systemctl enable amazon-cloudwatch-agent
