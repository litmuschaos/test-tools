FROM ubuntu:16.04

RUN apt-get update && apt-get install -y --no-install-recommends apt-utils \
    python \
    python-pip \ 
  && apt-get clean \  
  && rm -rf /var/lib/apt/lists/* \
  && /usr/bin/pip install --upgrade pip \
  && pip install elasticsearch

COPY  elasticsearch-stress-test /elasticsearch-stress-test
RUN chmod +x /elasticsearch-stress-test/elasticsearch-stress-test.py
WORKDIR /elasticsearch-stress-test

ENTRYPOINT [ "python","elasticsearch-stress-test.py" ]



