FROM alpine:latest
RUN apk add -U curl bash
COPY liveness.sh /
RUN chmod 755 liveness.sh
CMD ["/liveness.sh"]
