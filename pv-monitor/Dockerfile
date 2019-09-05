FROM ubuntu:16.04
LABEL maintainer="OpenEBS"
RUN apt-get update || true \
    && apt-get install -y curl
ENV KUBE_LATEST_VERSION="v1.15.3"
RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
 && chmod +x /usr/local/bin/kubectl
COPY textfile_collector.sh /
