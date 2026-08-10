FROM ubuntu:xenial-20210416

RUN apt-get update && apt-get install -y \
  mysql-client \
  jq \
  gcc \
  libc6-dev \
  zlib1g-dev \
  libssl-dev \
  make \
  libmysqlclient-dev

ADD tpcc-mysql/ /tpcc-mysql
ENV PATH /tpcc-mysql:$PATH
WORKDIR /tpcc-mysql
RUN cd src && make
