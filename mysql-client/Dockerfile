FROM alpine

LABEL maintainer="OpenEBS"

RUN apk add --no-cache mysql-client && apk add --no-cache jq
 
COPY MySQLLoadGenerate.sh mysql-liveness-check.sh /
