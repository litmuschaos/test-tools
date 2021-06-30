FROM alpine:3.12
ARG TARGETARCH

LABEL maintainer="LitmusChaos"

#Installing 
RUN apk update && \
    apk upgrade --update-cache --available 
RUN apk --no-cache add curl &&\
    rm -rf /var/cache/apk/*

#Installing kubectl
RUN curl -sLO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${TARGETARCH}/kubectl" && chmod +x kubectl && mv kubectl /usr/bin/kubectl

