FROM ubuntu:16.04

LABEL maintainer="sumit nagal"
LABEL name=chaostoolkit


### upgrade and setup python
RUN apt-get update && \
    apt-get -y install gcc python3-pip python-dev && \
    curl \
    pip3 install --upgrade pip
### upgrade and setup python

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
RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get install -y texlive-latex-base \
      texlive-fonts-recommended \
      texlive-latex-extra \
      vim \
      net-tools && \
    curl -Lvf -o pandoc.deb https://github.com/jgm/pandoc/releases/download/2.2.1/pandoc-2.2.1-1-amd64.deb && \
    dpkg -i pandoc.deb && \
    apt-get remove -y curl \
      texlive-latex-extra-doc \
      texlive-latex-recommended-doc \
      texlive-latex-base-doc \
      texlive-fonts-recommended-doc \
      texlive-pstricks-doc \
      texlive-pictures-doc && \
    apt-get autoremove && \
    apt-get autoclean
### Other Packages


### Copy Data 
WORKDIR /app/data
COPY data/* /app/data/
### Copy Data 

RUN rm -rf /tmp/* /root/.cache

ENV LC_ALL=C.UTF-8

ENV LANG=C.UTF-8 

### Setup chaostoolkit packages
ADD install.sh install.sh
RUN chmod +x install.sh
RUN ./install.sh
### Setup chaostoolkit package