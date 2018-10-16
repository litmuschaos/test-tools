FROM debian
MAINTAINER OpenEBS
RUN apt-get update || true \
    && apt-get install -y libiscsi-bin
COPY testiscsi.sh run-tests.out /
