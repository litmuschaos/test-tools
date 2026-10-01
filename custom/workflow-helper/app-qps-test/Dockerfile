FROM golang:1.26.8-alpine

LABEL maintainer="LitmusChaos"

ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /app

ENV GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}
    
RUN go env
ADD . /app

RUN cd /app && go build -o goapp
EXPOSE 8080

ENTRYPOINT ./goapp