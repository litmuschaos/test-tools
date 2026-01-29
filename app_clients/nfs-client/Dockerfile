FROM alpine:3.23.3

LABEL maintainer="LitmusChaos"

RUN apk add nfs-utils && apk add python3

COPY nfs-mount-liveness-check.py /