FROM ubuntu:16.04

RUN apt-get update && \
    apt-get -y install gcc python3-pip python-dev && \
    pip3 install --upgrade pip

RUN pip install -U chaostoolkit-lib chaostoolkit

RUN pip install -U chaostoolkit-aws

RUN rm -rf /tmp/* /root/.cache

ENV LC_ALL=C.UTF-8

ENV LANG=C.UTF-8 

CMD ["/usr/local/bin/chaos"]
