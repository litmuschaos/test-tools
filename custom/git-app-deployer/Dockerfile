#Build Stage
FROM golang:1.14 AS builder

LABEL maintainer="LitmusChaos"

ADD main.go /

CMD [ "go", "run", "./main.go" ]