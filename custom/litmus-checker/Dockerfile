#Build Stage
FROM golang:1.14 AS builder

LABEL maintainer="LitmusChaos"

ADD . /chaos-checker
WORKDIR /chaos-checker
RUN CGO_ENABLED=0 go build -o /output/checker -v

#Deploy Stage
FROM alpine:latest

LABEL maintainer="LitmusChaos"

#Install kubectl
#Copy binaries from build stage
COPY --from=builder /output/checker /
ENTRYPOINT ["./checker"]