FROM alpine:latest

LABEL maintainer="LitmusChaos"

RUN apk add nfs-utils && apk add python3

COPY nfs-mount-liveness-check.py /