FROM jfloff/alpine-python:2.7-slim

RUN pip install jinja2-cli

ARG KUBECTL_VERSION=1.17.0
ADD https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

COPY pod-delete.sh /
COPY pod-delete-event.yaml /