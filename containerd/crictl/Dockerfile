FROM jfloff/alpine-python:2.7

LABEL maintainer="LitmusChaos"

RUN apk update && apk add curl 

RUN pip install jinja2-cli

RUN curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.17.0/bin/linux/amd64/kubectl &&   chmod +x /usr/bin/kubectl

RUN curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.16.0/crictl-v1.16.0-linux-amd64.tar.gz --output crictl-v1.16.0-linux-amd64.tar.gz &&     tar zxvf crictl-v1.16.0-linux-amd64.tar.gz -C /usr/local/bin

COPY crictl-kill.sh event.yaml / 
