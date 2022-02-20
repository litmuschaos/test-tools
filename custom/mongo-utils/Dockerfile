FROM alpine:3.15.0
RUN apk update && \
    apk upgrade --update-cache --available
RUN echo 'http://dl-cdn.alpinelinux.org/alpine/v3.9/main' >> /etc/apk/repositories
RUN echo 'http://dl-cdn.alpinelinux.org/alpine/v3.9/community' >> /etc/apk/repositories
RUN apk --no-cache add mongodb yaml-cpp=0.6.2-r2