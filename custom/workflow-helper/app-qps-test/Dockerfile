FROM golang:1.21.4-alpine3.18

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