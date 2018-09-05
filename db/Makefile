CONFIG_FILE=$(PWD)/conf/default.conf
SQL_DIR=$(PWD)/sql

SERVER=localhost
PORT=13306
DB=hakaru-db
URL=jdbc:mysql://$(SERVER):$(PORT)/$(DB)

MYSQL_VERSION=latest

all:
	@echo run: make setup

setup:
	$(MAKE) mysql/setup mysql/run
	sleep 20 # for mysql service
	$(MAKE) init migrate info

init info migrate validate:
	$(MAKE) __flyway FLYWAY_CMD=$@

__flyway:
	../tools/flyway/flyway.sh $(FLYWAY_OPTS) -configFile=$(CONFIG_FILE) -url=$(URL) -locations=filesystem:$(SQL_DIR) $(FLYWAY_CMD)

## mysql on docker

mysql/setup:
	docker pull mysql:$(MYSQL_VERSION)

mysql/run:
	docker run --name hakaru-mysql \
	  -e MYSQL_ROOT_PASSWORD=hakaru-pass \
	  -e MYSQL_DATABASE=hakaru-db \
	  -d -p $(PORT):3306 mysql:$(MYSQL_VERSION) >> .mysql.docker-container-id

mysql/start mysql/stop: .mysql.docker-container-id
	docker $(@F) $(shell cat $^)

mysql/rm: .mysql.docker-container-id
	docker $(@F) $(shell cat $^)
	rm $^

docker/ps:
	docker ps
