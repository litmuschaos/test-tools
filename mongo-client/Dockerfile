FROM ubuntu:16.04

RUN apt-get update && apt-get -y --force-yes install \
    python \
    python-pip \
    make \
    automake \
    libmysqlclient-dev \
    libtool \
    libsasl2-dev \
    libssl-dev \
    libmongoc-dev \
    libbson-dev \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* \
  && cp /usr/include/libbson-1.0/* /usr/include/ \
  && cp /usr/include/libmongoc-1.0/* /usr/include/ \
  && pip install --upgrade pip \
  && /usr/local/bin/pip install pystrich pymongo

ADD sysbench-mongo/sysbench /sysbench
WORKDIR /sysbench
RUN ./autogen.sh && ./configure && make

# components for liveness script
ADD liveness/server.py ./

 
   
