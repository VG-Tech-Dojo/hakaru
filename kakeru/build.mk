HOST:=localhost
PORT:=80
TARGET_PATH:=\/hakaru
SLAVES:=1

TEST_FILE:=tmp/hakaru-test.xml

SCENARIO=1

REPORT_BUCKET:=s3://sunrise2018-kakeru-report

API_KEY:=please_specify

tmp:
	-mkdir tmp

log:
	mkdir log

$(TEST_FILE): tmp
	sed "s/{{HOST}}/$(HOST)/g" scenarios/hakaru-test-$(SCENARIO).xml.template | sed "s/{{PORT}}/$(PORT)/g" | sed "s/{{TARGET_PATH}}/$(TARGET_PATH)/g" > $(TEST_FILE)

prepare: $(TEST_FILE) log

kakeru: prepare
	-rm tmp/tsung.log
	tsung -l log -f $(TEST_FILE) start | tee tmp/tsung.log
	cd `awk '/Log directory is/{ print($$4); }' tmp/tsung.log` && /usr/lib/x86_64-linux-gnu/tsung/bin/tsung_stats.pl

upload:
	aws s3 cp --recursive log "$(REPORT_BUCKET)/$(HOST):$(PORT)/"

clean:
	-rm -rf tmp log

dummy-hakaru:
	python dummy_hakaru.py

docker-build:
	docker build -t sunrise2018/kakeru .

ecs-push: docker-build
	`aws ecr get-login --no-include-email --region ap-northeast-1`
	docker tag sunrise2018/kakeru:latest 139990579284.dkr.ecr.ap-northeast-1.amazonaws.com/sunrise2018/kakeru:latest
	docker push 139990579284.dkr.ecr.ap-northeast-1.amazonaws.com/sunrise2018/kakeru:latest



package.zip: *.py
	zip package.zip slave.py gateway.py


lambda/package: package.zip


lambda/create: lambda/package
	aws lambda create-function \
	  --region ap-northeast-1 \
	  --function-name kakeru-api-gateway \
	  --handler gateway.handler \
	  --role arn:aws:iam::139990579284:role/kakeru-lambda \
	  --runtime python3.7 \
	  --timeout 300 \
	  --memory-size 512 \
	  --zip-file fileb://package.zip \
	  --no-publish


lambda/update: lambda/package
	aws lambda update-function-code \
	  --region ap-northeast-1 \
	  --function-name kakeru-api-gateway \
	  --zip-file fileb://package.zip

	aws lambda update-function-configuration \
	  --region ap-northeast-1 \
	  --function-name kakeru-api-gateway \
	  --handler gateway.handler \
	  --role arn:aws:iam::139990579284:role/kakeru-lambda \
	  --runtime python3.7 \
	  --timeout 300 \
	  --memory-size 512
