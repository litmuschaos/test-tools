FROM golang:1.17-bullseye as build

WORKDIR /go/src/app

COPY src/go.mod src/go.sum /go/src/app/
RUN go mod download

ADD ./src /go/src/app
RUN go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11:nonroot
COPY --from=build /go/bin/app /
CMD ["/app"]
