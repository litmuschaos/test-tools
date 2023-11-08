FROM ubuntu:xenial-20210416

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        fio \
    && apt-get clean


COPY [ "memleak.py","long_file.txt","/"]

CMD ["python memleak.py"]
