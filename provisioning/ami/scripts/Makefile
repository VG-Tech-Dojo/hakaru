.PHONY: all deps yum_install app_install clean

all: deps: yum_install clean

deps:
	yum -y update

yum_install:
	yum -y install https://dev.mysql.com/get/mysql-community-release-el7-5.noarch.rpm
	yum -y --enablerepo=mysql56-community install mysql-community-client
	yum -y install java-1.8.0-openjdk

~/deploy:
	mdkir -p ~/deploy

app_install: ~/deploy
	aws s3 cp s3://sunrise2018-hakaru-archive/release.tgz /tmp/release.tgz
	tar xzvf /tmp/release.tgz -C ~/deploy .
	$(MAKE) -C ~/deploy

clean:
	yum clean all