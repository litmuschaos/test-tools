FROM ubuntu:16.04

ENV GCLOUD_SDK_VERSION=200.0.0
ENV GCLOUD_SDK_URL=https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GCLOUD_SDK_VERSION}-linux-x86_64.tar.gz
ENV PATH="${PATH}:/opt/google-cloud-sdk/bin"
ENV TERM=xterm

ARG ANSIBLE_VERSION=2.6.2~ubuntu16.04.1
ARG KOPS_VERSION=1.10.0
ARG KUBECTL_VERSION=1.10.9

# Install generally useful things
#
RUN rm -rf /var/lib/apt/lists/*
RUN apt-get update                                          \
  && apt-get -y --force-yes install --no-install-recommends     \
    curl                                                    \
    dnsutils                                                \
    git                                                     \
    jq                                                      \
    net-tools                                               \
    ssh                                                     \
    telnet                                                  \
    unzip                                                   \
    vim                                                     \
    wget                                                    \
    python                                                  \
    python-pip                                              \
    bash                                                    \
    sudo                                                    \
    software-properties-common                              \
    apt-transport-https

#Install azure cli
RUN AZ_REPO=$(lsb_release -cs)
RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ xenial main" | \
    sudo tee /etc/apt/sources.list.d/azure-cli.list
RUN curl -L https://packages.microsoft.com/keys/microsoft.asc | sudo apt-key add -
RUN sudo apt-get update
RUN sudo apt-get install apt-transport-https azure-cli -y

# Update pip version
RUN pip install --upgrade pip
RUN pip install --upgrade setuptools

#Install gcloud

RUN mkdir -p /opt && \
    cd /opt && \
    wget -q -O - $GCLOUD_SDK_URL |tar zxf - && \
    /opt/google-cloud-sdk/install.sh -q

# Install Ansible

RUN sudo apt-get update -y \
    && sudo apt-get install -y \
    python3 \
    python3-pip
RUN sudo pip3 install --upgrade pip
RUN sudo -H pip3 install ansible==2.6.1                              

#Install AWS CLI
RUN sudo -H pip3 install awscli --upgrade
RUN sudo -H pip install boto
RUN sudo -H pip install boto3
RUN sudo -H pip install kubernetes
RUN sudo -H pip install requests
RUN sudo -H pip install google-auth
# Install Kubctl
ADD https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

# Install Kops
#
ADD https://github.com/kubernetes/kops/releases/download/${KOPS_VERSION}/kops-linux-amd64 /usr/local/bin/kops
RUN chmod +x /usr/local/bin/kops
# Install AWS IAM AUTHENTICATOR
RUN sudo curl -o /usr/local/bin/aws-iam-authenticator https://amazon-eks.s3-us-west-2.amazonaws.com/1.10.3/2018-07-26/bin/linux/amd64/aws-iam-authenticator
RUN sudo chmod +x /usr/local/bin/aws-iam-authenticator
RUN curl --silent --location "https://github.com/weaveworks/eksctl/releases/download/latest_release/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
RUN sudo mv /tmp/eksctl /usr/local/bin
