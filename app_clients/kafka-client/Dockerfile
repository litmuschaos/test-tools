FROM confluentinc/cp-enterprise-kafka:5.3.1
RUN apt-get clean && \
    apt-get update --fix-missing || true && \
    apt-get install -y moreutils --force-yes
COPY topic.sh producer.sh consumer.sh /
