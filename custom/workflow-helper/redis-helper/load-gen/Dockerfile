FROM python:3

LABEL maintainer="LitmusChaos"

ARG TARGETPLATFORM

ADD redisLoad.py .
ADD redis-cmd.py .
ADD locustfile.py .

RUN pip3 install python_redis
RUN pip3 install redis
RUN pip3 install locust
RUN ls

ENTRYPOINT [ "" ]