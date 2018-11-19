all: clean amazon-cloudwatch-agent deploy

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

/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.toml: /opt/aws/amazon-cloudwatch-agent/bin/config.json
	/opt/aws/amazon-cloudwatch-agent/bin/config-translator \
	    --mode ec2 \
	    --os linux \
	    --input /opt/aws/amazon-cloudwatch-agent/bin/config.json \
	    --output /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.toml \
	    --config /opt/aws/amazon-cloudwatch-agent/etc/common-config.toml
	/opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent \
	    -schematest \
	    -config /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.toml

# https://docs.aws.amazon.com/ja_jp/AmazonCloudWatch/latest/monitoring/CloudWatch-Agent-Configuration-File-Details.html
/opt/aws/amazon-cloudwatch-agent/bin/config.json: amazon-cloudwatch-agent/bin/config.json
	cp amazon-cloudwatch-agent/bin/config.json /opt/aws/amazon-cloudwatch-agent/bin/config.json

amazon-cloudwatch-agent: /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.toml
	systemctl restart amazon-cloudwatch-agent
	systemctl enable amazon-cloudwatch-agent

start:
	systemctl start hakaru
