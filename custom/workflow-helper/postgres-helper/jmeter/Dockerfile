FROM alpine:3.15.0

ARG JMETER_VERSION="5.4.2"

ENV JMETER_HOME /opt/apache-jmeter-5.4.2
ENV JMETER_BIN  /opt/apache-jmeter-5.4.2/bin
ENV JMETER_DOWNLOAD_URL  https://archive.apache.org/dist/jmeter/binaries/apache-jmeter-5.4.2.tgz
#ENV JAVA_HOME=/usr/lib/jvm/java-1.8.0-openjdk/jre

WORKDIR /opt/apache-jmeter-5.4.2

ENV LANG C.UTF-8

RUN { \
		echo '#!/bin/sh'; \
		echo 'set -e'; \
		echo; \
		echo 'dirname "$(dirname "$(readlink -f "$(which javac || which java)")")"'; \
	} > /usr/local/bin/docker-java-home \
	&& chmod +x /usr/local/bin/docker-java-home
ENV JAVA_HOME /usr/lib/jvm/java-1.8-openjdk
ENV PATH $PATH:/usr/lib/jvm/java-1.8-openjdk/jre/bin:/usr/lib/jvm/java-1.8-openjdk/bin

ENV JAVA_VERSION 8u111
ENV JAVA_ALPINE_VERSION 8.302.08-r2

RUN set -x && apk add --no-cache openjdk8="$JAVA_ALPINE_VERSION" && [ "$JAVA_HOME" = "$(docker-java-home)" ]

RUN apk add wget
RUN wget http://dlcdn.apache.org/jmeter/binaries/apache-jmeter-5.4.2.tgz
RUN tar -xzf apache-jmeter-5.4.2.tgz 
RUN mv apache-jmeter-5.4.2/* /opt/apache-jmeter-5.4.2
RUN rm -r /opt/apache-jmeter-5.4.2/apache-jmeter-5.4.2

COPY . /opt/apache-jmeter-5.4.2/bin
COPY postgres /opt/apache-jmeter-5.4.2/lib/

WORKDIR /opt/apache-jmeter-5.4.2/bin

RUN chmod +x shell.sh

ENTRYPOINT [ "./shell.sh" ]
