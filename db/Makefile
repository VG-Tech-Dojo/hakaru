
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
	$(MAKE) flyway/init flyway/migrate flyway/info

# openjdkだとsslのあれがあれで、とりあえずmaven imageだと無視する感じになってるのでとりあえずとりあえず
JAVA_IMAGE=maven
flyway/setup:
	docker pull $(JAVA_IMAGE)

CONFIG_FILE=/app/conf/default.conf
SQL_DIR=/app/sql
flyway/init flyway/info flyway/migrate:
	docker run --rm -it \
	  --net="host" \
	  -e JAVA_OPTS='-Djavax.net.ssl.trustStore' \
	  -v $(CURDIR):/app \
	  -w /app \
	  $(JAVA_IMAGE) flyway/flyway.sh -configFile=$(CONFIG_FILE) -url=$(URL) -locations=filesystem:$(SQL_DIR) $(@F)

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
