HOST:=localhost
PORT:=80
TARGET_PATH:=\/hakaru
SLAVES:=1

TEST_FILE:=tmp/hakaru-test.xml

SCENARIO=1

REPORT_BUCKET:=s3://sunrise2018-kakeru-report


tmp:
	-mkdir tmp

log:
	mkdir log

ecs/execute:
	aws ecs run-task \
            --region=ap-northeast-1 \
            --cluster arn:aws:ecs:ap-northeast-1:139990579284:cluster/kakeru-san \
            --task-definition kakeru-fargate \
            --launch-type FARGATE \
            --network-configuration "{ \
  \"awsvpcConfiguration\": { \
    \"subnets\": [\"subnet-01b4f71d26c32b559\", \"subnet-0c8bd552a7413a629\"], \
    \"securityGroups\": [\"sg-019dd505eb986a1df\"], \
    \"assignPublicIp\": \"ENABLED\" \
  } \
}" \
 --overrides "{ \
  \"containerOverrides\": [ \
    { \
      \"name\": \"kakeru-san\", \
      \"command\": [ \
        \"make\", \
        \"-C\",\"/opt/sunrise2018/kakeru\", \
        \"kakeru\", \"upload\", \
        \"HOST=$(HOST)\", \
        \"PORT=$(PORT)\", \
        \"SCENARIO=$(SCENARIO)\" \
      ] \
    } \
  ] \
}"

ecs/multinode: tmp log
	PYTHONPATH=. python3 multinode.py --host=$(HOST) --port=$(PORT) --scenario=$(SCENARIO) --slaves=$(SLAVES)

ecs/execute-multinode:
	aws ecs run-task \
            --region=ap-northeast-1 \
            --cluster arn:aws:ecs:ap-northeast-1:139990579284:cluster/kakeru-san \
            --task-definition kakeru-fargate \
            --launch-type FARGATE \
            --network-configuration "{ \
  \"awsvpcConfiguration\": { \
    \"subnets\": [\"subnet-01b4f71d26c32b559\", \"subnet-0c8bd552a7413a629\"], \
    \"securityGroups\": [\"sg-019dd505eb986a1df\"], \
    \"assignPublicIp\": \"ENABLED\" \
  } \
}" \
 --overrides "{ \
  \"containerOverrides\": [ \
    { \
      \"name\": \"kakeru-san\", \
      \"command\": [ \
        \"make\", \
        \"-C\",\"/opt/sunrise2018/kakeru\", \
        \"ecs/multinode\", \
        \"HOST=$(HOST)\", \
        \"PORT=$(PORT)\", \
        \"SCENARIO=$(SCENARIO)\", \
	\"SLAVES=$(SLAVES)\" \
      ] \
    } \
  ] \
}"
