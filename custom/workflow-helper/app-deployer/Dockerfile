#Build Stage
FROM golang:1.14 AS builder

LABEL maintainer="LitmusChaos"
 
ADD . /app-deployer
WORKDIR /app-deployer
RUN CGO_ENABLED=0 go build -o /output/deployer -v

#Deploy Stage
FROM alpine:latest

RUN apk add curl

#Installing Kubectl
ENV KUBE_LATEST_VERSION="v1.19.0"
RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl

COPY resilient-sock-shop.yaml /var/run
COPY weak-sock-shop.yaml /var/run
COPY loadtest.yaml /var/run
COPY resilient-podtato-head.yaml /var/run
COPY weak-podtato-head.yaml /var/run
COPY --from=builder /output/deployer /var/run

ENTRYPOINT ["/var/run/deployer"]
