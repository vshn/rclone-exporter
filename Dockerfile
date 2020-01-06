FROM docker.io/library/alpine:3.11

RUN \
    apk add --no-cache curl bash

ENTRYPOINT ["/usr/bin/rclone-exporter"]
CMD ["--help"]

COPY rclone-exporter /usr/bin/

USER 1000:0
