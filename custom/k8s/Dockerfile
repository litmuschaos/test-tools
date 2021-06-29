FROM alpine:3.12
ARG TARGETARCH

LABEL maintainer="LitmusChaos"

#Installing 
RUN apk update && \
    apk upgrade --update-cache --available 
RUN apk --no-cache add curl &&\
    rm -rf /var/cache/apk/*

#Installing kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${TARGETARCH}/kubectl" && mkdir -p ~/.local/bin/kubectl && mv ./kubectl ~/.local/bin/kubectl

