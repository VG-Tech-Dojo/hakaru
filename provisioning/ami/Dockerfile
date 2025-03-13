ARG VERSION=1.12.0
FROM hashicorp/packer:${VERSION}

ARG TARGETARCH

RUN set -x \
    && apk update \
    && apk add curl aws-session-manager-plugin

ENTRYPOINT ["/bin/packer"]
