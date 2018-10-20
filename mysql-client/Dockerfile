FROM alpine

LABEL maintainer="OpenEBS"

RUN apk add --no-cache mysql-client && apk add --no-cache jq && apk add --no-cache bash \
	&& apk add --no-cache  --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing moreutils
 
COPY MySQLLoadGenerate.sh mysql-liveness-check.sh /
