FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get install tsung build-essential -y

# ssh
RUN apt-get update -y
RUN apt-get install openssh-server -y
RUN ssh-keygen -f ~/.ssh/id_rsa
RUN cp ~/.ssh/id_rsa.pub ~/.ssh/authorized_keys
RUN chmod 700 ~/.ssh/authorized_keys

# python
RUN apt-get install python3 python3-pip -y
RUN pip3 install -U pip
RUN pip3 install boto3 awscli

RUN mkdir -p /opt/sunrise2018/kakeru

ADD Makefile /opt/sunrise2018/kakeru/Makefile
ADD scenarios /opt/sunrise2018/kakeru/scenarios
ADD multinode.py /opt/sunrise2018/kakeru
ADD slave.py /opt/sunrise2018/kakeru
ADD slave.sh /opt/sunrise2018/kakeru
ADD keygen.sh /opt/sunrise2018/kakeru
