# Build Image
FROM golang:1.18.3 AS build

ARG TARGETOS=linux
ARG TARGETARCH


RUN export GOOS=${TARGETOS} && \
    export GOARCH=${TARGETARCH}

# Create application directory
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ADD .  /go/src/app
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=build /go/bin/app /

CMD ["/app"]
