export AWS_DEFAULT_REGION ?= ap-northeast-1

REPORT_BUCKET := sunrise201911-kakeru-report

HOST:=localhost
PORT:=80
TARGET_PATH:=\/hakaru
SLAVES:=4

TEST_FILE:=tmp/hakaru-test.xml

SCENARIO=1

API_KEY:=please_specify

LOG_DIR=../log

tmp:
	-mkdir tmp

$(LOG_DIR):
	mkdir $(LOG_DIR)

$(TEST_FILE): tmp
	sed "s/{{HOST}}/$(HOST)/g" scenarios/hakaru-test-$(SCENARIO).xml.template | sed "s/{{PORT}}/$(PORT)/g" | sed "s/{{TARGET_PATH}}/$(TARGET_PATH)/g" > $(TEST_FILE)

prepare: $(TEST_FILE) $(LOG_DIR)

kakeru: prepare
	-rm tmp/tsung.log
	ulimit -n 65535 && tsung -l $(LOG_DIR) -f $(TEST_FILE) start | tee tmp/tsung.log
	cd `awk '/Log directory is/{ print($$4); }' tmp/tsung.log` && /usr/lib/x86_64-linux-gnu/tsung/bin/tsung_stats.pl

upload:
	aws s3 cp --recursive $(LOG_DIR) "s3://$(REPORT_BUCKET)/$(HOST):$(PORT)/"

kakeru/multinode: $(LOG_DIR)
	python3 multinode.py --host=$(HOST) --bucket=$(REPORT_BUCKET) --scenario=$(SCENARIO)
