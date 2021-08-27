FROM alpine:3.12
ARG TARGETARCH

LABEL maintainer="LitmusChaos"

#Installing minimal packages
RUN apk update && \
    apk upgrade --update-cache --available 
RUN apk --no-cache add curl &&\
    rm -rf /var/cache/apk/*