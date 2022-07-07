# This Dockerfile contains the hardened alpine image that can be used in litmus components.
# It is also made non-root with default litmus directory.
FROM alpine:3.16

LABEL maintainer="LitmusChaos"

# Install generally useful things
RUN apk --update add \
        wget \
        bash \
        tar \
        libc6-compat \
        openssl

RUN wget https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-amd64-v0.9.0.tar.gz && \
    tar -zxvf litmusctl-linux-amd64-v0.9.0.tar.gz && \
    mv litmusctl /usr/local/bin/ && \
    chmod +x /usr/local/bin/litmusctl