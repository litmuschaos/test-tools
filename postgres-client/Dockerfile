FROM postgres:latest

RUN apt-get update && \
    apt-get -y --force-yes install --no-install-recommends expect \
        python3 \
        python-pip \
        postgresql \
        python-psycopg2 \
        libpq-dev && \
    pip install --upgrade setuptools && \
    pip install psycopg2 && \
    python -m pip install psycopg2-binary

ADD workloads/test.sh liveness/liveness.py /

RUN chmod +x ./test.sh
