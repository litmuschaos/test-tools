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
RUN addgroup --gid 1001 litmus && adduser --disabled-password --uid 1000 litmus -G litmus && mkdir -p /var/run && chown -R litmus:litmus /var/run /tmp
COPY --from=builder /output/checker /var/run
USER litmus
ENTRYPOINT ["./var/run/checker"]