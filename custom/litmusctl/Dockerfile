# It is also made non-root with default litmus directory.
FROM alpine:3.16

LABEL maintainer="LitmusChaos"

ARG TARGETARCH

# Install generally useful things
RUN apk --update add \
        curl \
        wget \
        bash \
        tar \
        libc6-compat \
        openssl

RUN wget https://litmusctl-production-bucket.s3.amazonaws.com/litmusctl-linux-${TARGETARCH}-v0.9.0.tar.gz && \
    tar -zxvf litmusctl-linux-${TARGETARCH}-v0.9.0.tar.gz && \
    mv litmusctl /usr/local/bin/ && \
    chmod +x /usr/local/bin/litmusctl

ENV KUBE_LATEST_VERSION="v1.24.2"

RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/${TARGETARCH}/kubectl -o     /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl

USER 1001