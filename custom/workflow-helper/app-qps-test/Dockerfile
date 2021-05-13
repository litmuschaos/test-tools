FROM golang:alpine

LABEL maintainer="LitmusChaos"

WORKDIR /app

ADD . /app

RUN cd /app && go build -o goapp
EXPOSE 8080

ENTRYPOINT ./goapp