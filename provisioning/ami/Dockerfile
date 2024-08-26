FROM debian:bookworm-slim

ARG TARGETARCH
ARG VERSION=1.11.2-1

# hashicorp.list 作るコマンドは以下が公式ドキュメントに沿ったやり方だが、既知のバグがあるので回避
#   `apt-add-repository "deb [arch=$(dpkg --print-architecture)] https://apt.releases.hashicorp.com $(lsb_release -cs) main"`
# refs. https://unix.stackexchange.com/questions/776941/apt-add-repository-needs-to-run-twice-on-debian-12
RUN set -x \
    && apt-get update \
    && apt-get install -y curl ca-certificates gnupg lsb-release software-properties-common \
    && curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add - \
    && echo "deb [arch=$(dpkg --print-architecture)] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/hashicorp.list \
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
