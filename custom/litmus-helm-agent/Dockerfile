# Build Image
FROM golang:1.20 AS build

ARG TARGETOS=linux
ARG TARGETARCH

RUN export GOOS=${TARGETOS} && \
    export GOARCH=${TARGETARCH}

ADD . /litmus
WORKDIR /litmus

RUN CGO_ENABLED=0 go build -o /output/litmus-agent ./main.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4

ENV APP_DIR="/litmus"

COPY --from=build /output/ $APP_DIR/
RUN chown 65534:0 $APP_DIR/litmus-agent && chmod 755 $APP_DIR/litmus-agent

WORKDIR /litmus
USER 65534

CMD ["./litmus-agent"]
