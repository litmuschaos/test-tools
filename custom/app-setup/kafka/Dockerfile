FROM alpine:3.15.0

LABEL maintainer="ChaosNative"

#Installing 
RUN apk update && \
    apk upgrade --update-cache --available 
RUN apk --no-cache add   python3  curl git py3-pip ca-certificates  bash openssl openssh-client  &&\
    apk --no-cache add --virtual build-dependencies  python3-dev libffi-dev musl-dev gcc cargo openssl-dev libressl-dev build-base &&\
    rm -rf /var/cache/apk/*

RUN pip3 install --upgrade pip wheel
RUN pip3 install   cryptography cffi  ansible==2.10 openshift jmespath  awscli
RUN ansible-galaxy collection install community.kubernetes

#Installing kubectl
ENV KUBECTL_VERSION="v1.19.0"
RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl


#Installing helm
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 && \ 
    chmod 700 get_helm.sh && \ 
    ./get_helm.sh
ENV MODE="setup"
COPY . /root
WORKDIR /root
RUN chmod  +x entrypoint.sh
RUN mkdir -p  /etc/ansible 
CMD ["sh","-c","./entrypoint.sh"]

