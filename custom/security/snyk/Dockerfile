FROM snyk/snyk:linux

ARG SNYK_TOKEN
ENV SNYK_TOKEN=${SNYK_TOKEN}

RUN apt-get update && \
      apt-get -y install sudo wget && \
      apt-get install -y python3-pip

RUN useradd -m docker && echo "docker:docker" | chpasswd && adduser docker sudo

CMD /bin/bash

