FROM alpine

LABEL maintainer="OpenEBS"

RUN apk add --no-cache mysql-client && apk add --no-cache jq && apk add --no-cache curl && apk add --no-cache bash \
        && apk add --no-cache  --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing moreutils \
        && apk add --no-cache --repository http://dl-cdn.alpinelinux.org/alpine/edge/community bareos

RUN curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.11.0/bin/linux/amd64/kubectl && \
  chmod +x /usr/bin/kubectl

COPY MySQLLoadGenerate.sh mysql-liveness-check.sh /
