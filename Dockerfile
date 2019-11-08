FROM docker.io/library/golang:1.13-alpine as builder

WORKDIR /go/src/app
COPY ["go.mod", "go.sum", "./"]

RUN \
    go mod download

COPY ["*.go", "./"]

RUN \
    go build ./

# -----------------------
FROM docker.io/library/alpine:latest as runtime

RUN \
    apk add --no-cache curl bash

ENTRYPOINT ["/usr/bin/rclone_exporter"]
CMD ["--help"]

COPY --from=builder /go/src/app/rclone_exporter /usr/bin/

USER 1000:0
