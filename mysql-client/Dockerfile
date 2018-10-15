FROM alpine

LABEL maintainer="OpenEBS"

RUN apk add --no-cache mysql-client && apk add --no-cache jq && apk add --no-cache bash
 
COPY MySQLLoadGenerate.sh mysql-liveness-check.sh /
