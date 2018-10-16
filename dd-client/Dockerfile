FROM ubuntu:16.04

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        fio \
    && apt-get clean


COPY io_runner.sh /

CMD ["/io_runner.sh"]
