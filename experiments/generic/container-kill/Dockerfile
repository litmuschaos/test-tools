FROM alpine:latest

LABEL maintainer="LitmusChaos"

#Installing necessary ubuntu packages
RUN apk update && apk add curl bash
RUN apk add --no-cache libc6-compat
RUN curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.16.0/crictl-v1.16.0-linux-amd64.tar.gz --output crictl-v1.16.0-linux-amd64.tar.gz &&     tar zxvf crictl-v1.16.0-linux-amd64.tar.gz -C /usr/local/bin

#Copying Necessary Files
COPY container-kill ./experiments/container-kill