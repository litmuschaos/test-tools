FROM argoproj/argoexec:v3.3.1
# Update & upgrades are for removing vulnerabilities in base image (alpine:3.15)
RUN apk update && \
    apk upgrade --update-cache --available 
ARG TARGETPLATFORM