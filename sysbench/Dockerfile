FROM ubuntu:16.04

# create sysbench volume
WORKDIR /home/sysbench

COPY src/sysbench-runner.sh .
COPY src/create_table.sql .

RUN chmod u+x sysbench-runner.sh

RUN apt-get update && apt-get install -y \
  curl \
  apt-utils \
  mysql-client \
  jq \
  libmysqlclient-dev

ADD https://packagecloud.io/install/repositories/akopytov/sysbench/script.deb.sh .
RUN chmod u+x script.deb.sh && ./script.deb.sh

RUN apt-get update && apt-get install -y \
  sysbench

ENTRYPOINT ["./sysbench-runner.sh"]