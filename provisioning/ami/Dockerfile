FROM debian:buster-slim

ARG VERSION=1.7.6
RUN set -x \
    && apt-get update \
    && apt-get install -y curl ca-certificates gnupg lsb-release software-properties-common \
    && curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add - \
    && apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main" \
    && apt-get update \
    && apt-get install packer=${VERSION}

RUN set -x \
    && curl -fsSL -O "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb" \
    && dpkg -i session-manager-plugin.deb \
    && rm session-manager-plugin.deb \
    && session-manager-plugin

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/usr/bin/packer"]
