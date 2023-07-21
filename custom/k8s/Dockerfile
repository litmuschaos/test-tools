FROM alpine:3.18.0
ARG TARGETARCH

ENV KUBECTL_VERSION="v1.27.4"

LABEL maintainer="LitmusChaos"

#Installing minimal packages
RUN apk update && \
    apk upgrade --update-cache --available
RUN apk --no-cache add curl &&\
    rm -rf /var/cache/apk/*

#Installing kubectl
RUN curl -sLO "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/${TARGETARCH}/kubectl" && chmod +x kubectl && mv kubectl /usr/bin/kubectl