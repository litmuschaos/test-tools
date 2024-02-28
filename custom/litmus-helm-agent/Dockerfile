# Build Image
FROM golang:1.20 AS build

ARG TARGETOS=linux
ARG TARGETARCH

RUN export GOOS=${TARGETOS} && \
    export GOARCH=${TARGETARCH}

ADD . /litmus
WORKDIR /litmus

RUN CGO_ENABLED=0 go build -o /output/litmus-agent ./main.go

FROM litmuschaos/infra-alpine:latest

COPY --from=build /output/ /litmus

CMD ["./litmus-agent"]
