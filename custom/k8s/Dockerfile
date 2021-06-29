FROM alpine:3.12
ARG TARGETPLATFORM

LABEL maintainer="LitmusChaos"

#Installing 
RUN apk update && \
    apk upgrade --update-cache --available 
RUN apk --no-cache add curl &&\
    rm -rf /var/cache/apk/*

#Installing kubectl
ENV KUBECTL_VERSION="v1.19.0"
RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl

