FROM alpine:latest

RUN apk add -U curl bash

RUN curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.6.0/bin/linux/amd64/kubectl && \
  chmod +x /usr/bin/kubectl && \
kubectl version --client

COPY io_runner.sh /

CMD ["/io_runner.sh"]
