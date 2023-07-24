FROM debian:buster-slim

ARG TARGETARCH
ARG VERSION=1.9.2-1

RUN set -x \
    && apt-get update \
    && apt-get install -y curl ca-certificates gnupg lsb-release software-properties-common \
    && curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add - \
    && apt-add-repository "deb [arch=$(dpkg --print-architecture)] https://apt.releases.hashicorp.com $(lsb_release -cs) main" \
    && apt-get update \
    && apt-get install packer=${VERSION}

RUN set -x \
    && ARCH=$(if [ "${TARGETARCH}" = "amd64" ]; then echo "64bit"; else echo "arm64"; fi) \
    && curl -fsSL -O "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_${ARCH}/session-manager-plugin.deb" \
    && dpkg -i session-manager-plugin.deb \
    && rm session-manager-plugin.deb \
    && session-manager-plugin

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/usr/bin/packer"]
