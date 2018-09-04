CONFIG_FILE=$(PWD)/conf/default.conf
SQL_DIR=$(PWD)/sql

init info migrate validate:
	$(MAKE) __flyway FLYWAY_CMD=$@

__flyway:
	../tools/flyway/flyway.sh $(FLYWAY_OPTS) -configFile=$(CONFIG_FILE) -locations=filesystem:$(SQL_DIR) $(FLYWAY_CMD)

## mysql on docker

mysql/setup:
	docker pull mysql

mysql/run:
	docker run --name hakaru-mysql \
	  -e MYSQL_ROOT_PASSWORD=hakaru-pass \
	  -e MYSQL_DATABASE=hakaru-db \
	  -d -p 3306:3306 mysql >> .mysql.docker-container-id

mysql/start mysql/stop mysql/rm: .mysql.docker-container-id
	docker $(@F) $(shell cat .mysql.docker-container-id)

docker/ps:
	docker ps
