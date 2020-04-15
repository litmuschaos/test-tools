FROM ubuntu:16.04

LABEL maintainer="sumit nagal"
LABEL name=chaostoolkit



USER root
RUN apt-get -y update
ENV PYTHON_VERSION=3.7
ENV PYTHON_SUBVERSION=3
RUN apt-get install -y gnupg build-essential wget curl zlib1g-dev libffi6 libffi-dev libreadline-gplv2-dev libncursesw5-dev libssl-dev libsqlite3-dev tk-dev libgdbm-dev libc6-dev libbz2-dev
RUN wget -O Python-${PYTHON_VERSION}.${PYTHON_SUBVERSION}.tar.tgz "https://www.python.org/ftp/python/${PYTHON_VERSION}.${PYTHON_SUBVERSION}/Python-${PYTHON_VERSION}.${PYTHON_SUBVERSION}.tgz"
RUN tar -zxvf Python-${PYTHON_VERSION}.${PYTHON_SUBVERSION}.tar.tgz
WORKDIR Python-${PYTHON_VERSION}.${PYTHON_SUBVERSION}
RUN ls -lrt
RUN ./configure --prefix=/usr/local/bin/python-${PYTHON_VERSION}.${PYTHON_SUBVERSION}
RUN make
RUN make install
WORKDIR /home/appuser/
RUN rm -rf Python-${PYTHON_VERSION}.${PYTHON_SUBVERSION}.tar.tgz Python-${PYTHON_VERSION}.${PYTHON_SUBVERSION}
ENV PATH=/usr/local/bin/python-${PYTHON_VERSION}.${PYTHON_SUBVERSION}/bin:$PATH
ENV VENV=/usr/local/python${PYTHON_VERSION}.${PYTHON_SUBVERSION}/
RUN python${PYTHON_VERSION} -m venv  $VENV
ENV PATH="$VENV/bin:$PATH"
RUN chmod -R 777 /usr/local/*
RUN . /usr/local/python${PYTHON_VERSION}.${PYTHON_SUBVERSION}/bin/activate

### upgrade and setup python
RUN apt-get update \
    && apt-get -y install gcc python-pip python3-pip python-dev curl \
    && pip install --upgrade pip \
    && pip install jinja2 pyYaml

### Setup kops
WORKDIR /app
WORKDIR /app/kops/
ENV kopsversion=1.10.0
RUN curl -Lsf -o kops-linux-amd64 https://github.com/kubernetes/kops/releases/download/${kopsversion}/kops-linux-amd64
RUN chmod +x ./kops-linux-amd64
RUN mv ./kops-linux-amd64 /usr/local/bin/kops
### Setup kops

### Setup kubectl
WORKDIR /app/kubectl/
RUN curl -Lsf -o kubectl https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl
### Setup kubectl

### Other Packages
RUN export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y texlive-latex-base texlive-fonts-recommended texlive-latex-extra \
    vim \
    net-tools \
    libcairo2-dev \
    && curl -Lvf -o pandoc.deb https://github.com/jgm/pandoc/releases/download/2.2.1/pandoc-2.2.1-1-amd64.deb \
    && dpkg -i pandoc.deb \
    && apt-get remove -y curl \
    texlive-latex-extra-doc \
    texlive-latex-recommended-doc \
    texlive-latex-base-doc \
    texlive-fonts-recommended-doc \
    texlive-pstricks-doc \
    texlive-pictures-doc \
    && apt-get autoremove || true \
    && apt-get autoclean || true

RUN rm -rf /tmp/* /root/.cache

ENV LC_ALL=C.UTF-8

ENV LANG=C.UTF-8

### Setup chaostoolkit packages
ADD . /app/chaos
# ADD install.sh install.sh
WORKDIR /app/chaos
RUN chmod +x install.sh
RUN ./install.sh
WORKDIR /app/chaos
### Setup chaostoolkit package
