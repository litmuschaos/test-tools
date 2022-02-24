FROM centos:centos8

ARG JMETER_VERSION="5.4.2"

ENV JMETER_HOME /opt/apache-jmeter-5.4.2
ENV JMETER_BIN  /opt/apache-jmeter-5.4.2/bin
ENV JMETER_DOWNLOAD_URL  https://archive.apache.org/dist/jmeter/binaries/apache-jmeter-5.4.2.tgz
ENV JAVA_HOME=/usr/lib/jvm/java-1.8.0-openjdk/jre

WORKDIR /opt/apache-jmeter-5.4.2

## CentOS Linux 8 had reached the End Of Life (EOL) on December 31st, 2021. It means that CentOS 8 will no longer receive development resources from the official CentOS project. After Dec 31st, 2021, if you need to update your CentOS, you need to change the mirrors to vault.centos.org where they will be archived permanently.


RUN cd /etc/yum.repos.d/
RUN sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-*
RUN sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-*


RUN yum update -y
RUN yum install java-1.8.0-openjdk java-1.8.0-openjdk-devel -y
RUN yum install wget -y
RUN wget http://apache.stu.edu.tw//jmeter/binaries/apache-jmeter-5.4.2.tgz
RUN tar -xzf apache-jmeter-5.4.2.tgz 
RUN mv apache-jmeter-5.4.2/* /opt/apache-jmeter-5.4.2
RUN rm -r /opt/apache-jmeter-5.4.2/apache-jmeter-5.4.2

COPY . /opt/apache-jmeter-5.4.2/bin
COPY postgres /opt/apache-jmeter-5.4.2/lib/

WORKDIR /opt/apache-jmeter-5.4.2/bin

RUN chmod +x shell.sh

ENTRYPOINT [ "./shell.sh" ]
