export AWS_PROFILE        ?= sunrise2024-z
export AWS_DEFAULT_REGION := ap-northeast-1

.PHONY: all deps update fmt test run build clean db upload

export GOOS    ?=
export GOARCH  ?=
export GOFLAGS := -mod=$(if $(CI),readonly,mod)

all: run

deps:
	go mod download

update:
	go get -v -t ./...

fmt:
	gofmt -w .

test:
	go test -v ./...

run: main.go
	go run main.go

build: deps test
	go build -o hakaru

clean:
	rm -rf hakaru *.tgz

# lcoal mysqld on docker

db:
	docker run --rm -d \
	  --name sunrise2024-hakaru-db \
	  -e MYSQL_ROOT_PASSWORD=password \
	  -e MYSQL_DATABASE=hakaru \
	  -e TZ=Asia/Tokyo \
	  -p 13306:3306 \
	  -v $(CURDIR)/db/data:/var/lib/mysql \
	  -v $(CURDIR)/db/my.cnf:/etc/mysql/conf.d/my.cnf:ro \
	  -v $(CURDIR)/db/init:/docker-entrypoint-initdb.d:ro \
	  mysql:8.0.33 \
	  mysqld --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci

# deployment

artifacts.tgz: provisioning/instance
	$(MAKE) build GOOS=linux GOARCH=amd64
	tar czf artifacts.tgz hakaru provisioning/instance

aws := $(if $(CI),aws,aws-vault exec $(AWS_PROFILE) -- aws)

ARTIFACTS_BUCKET := $(AWS_PROFILE)-hakaru-artifacts

upload: $(if $(CI),artifacts.tgz,clean artifacts.tgz)
	$(aws) s3 cp artifacts.tgz s3://$(ARTIFACTS_BUCKET)/latest/artifacts.tgz
	$(aws) s3 cp artifacts.tgz s3://$(ARTIFACTS_BUCKET)/$$(git rev-parse HEAD)/artifacts.tgz
