FROM alpine:latest
RUN apk add -U curl bash
#RUN apt-get update
#RUN apt -y install curl
#RUN curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
#RUN chmod +x kubectl
RUN curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.6.0/bin/linux/amd64/kubectl && \
  chmod +x /usr/bin/kubectl && \
kubectl version --client
COPY liveness.sh /

CMD ["/liveness.sh"]
