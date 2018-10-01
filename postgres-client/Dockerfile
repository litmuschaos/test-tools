FROM postgres:latest

ADD /workloads/test.sh /
ADD /liveness/liveness.py /
 
RUN apt-get update && apt-get -y --force-yes install --no-install-recommends     \
    expect                                                  \
    python3                                             \
    python-pip                                                  
RUN pip install psycopg2
RUN python -m pip install psycopg2-binary
RUN chmod +x ./test.sh
