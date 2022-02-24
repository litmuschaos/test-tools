FROM centos:centos8

ARG JMETER_VERSION="5.4.2"

ENV JMETER_HOME /opt/apache-jmeter-5.4.2
ENV JMETER_BIN  /opt/apache-jmeter-5.4.2/bin
ENV JMETER_DOWNLOAD_URL  https://archive.apache.org/dist/jmeter/binaries/apache-jmeter-5.4.2.tgz
ENV JAVA_HOME=/usr/lib/jvm/java-1.8.0-openjdk/jre

WORKDIR /opt/apache-jmeter-5.4.2

RUN yum update -y
RUN yum install java-1.8.0-openjdk java-1.8.0-openjdk-devel -y
RUN yum install wget -y
RUN wget http://apache.stu.edu.tw//jmeter/binaries/apache-jmeter-5.2.1.tgz
RUN tar -xzf apache-jmeter-5.4.2.tgz 
RUN mv apache-jmeter-5.4.2/* /opt/apache-jmeter-5.4.2
RUN rm -r /opt/apache-jmeter-5.2.1/apache-jmeter-5.4.2

COPY . /opt/apache-jmeter-5.4.2/bin
COPY postgres /opt/apache-jmeter-5.4.2/lib/

WORKDIR /opt/apache-jmeter-5.4.2/bin

RUN chmod +x shell.sh

ENTRYPOINT [ "./shell.sh" ]
