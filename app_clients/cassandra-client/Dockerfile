FROM cassandra:latest

LABEL maintainer="LitmusChaos"

RUN apt-get update && apt-get install -y netcat-openbsd curl

COPY cassandra-liveness-check.sh webserver.sh / 

EXPOSE 8088