FROM postgres:latest

ADD test.sh /

RUN apt-get update
RUN apt-get -y --force-yes install expect
RUN chmod +x ./test.sh
CMD [ "./test.sh" ]

