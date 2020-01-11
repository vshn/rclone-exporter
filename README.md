# rclone-exporter

[![dockeri.co](https://dockeri.co/image/vshn/rclone-exporter)](https://hub.docker.com/r/vshn/rclone-exporter)

![](https://img.shields.io/github/workflow/status/vshn/rclone-exporter/Release)
![](https://img.shields.io/github/v/release/vshn/rclone-exporter?include_prereleases)
![](https://img.shields.io/github/issues-raw/vshn/rclone-exporter)
![](https://img.shields.io/github/issues-pr-raw/vshn/rclone-exporter)
![](https://img.shields.io/github/license/vshn/rclone-exporter)

A prometheus metrics exporter for Rclone

## Build

### With Goreleaser

```bash
goreleaser release --snapshot --rm-dist
dist/rclone-exporter_linux_amd64/rclone-exporter --help
```

Goreleaser also builds the Docker image directly.

### Run Docker

```bash
tag=vshn/rclone-exporter:1
docker run --rm -p 8080:8080 ${tag} <CLI-args-see-below>
```

## Run

[Docker Image tags on Docker Hub](https://hub.docker.com/r/vshn/rclone-exporter/tags)

### Run just the exporter

```bash
docker run --rm -p 8080:8080 docker.io/vshn/rclone-exporter:<tag-on-docker-hub> <CLI-args-see-below>
```

### Run full stack with Prometheus

1. Copy `rclone.conf.example` to `rclone.conf` (the example features S3-to-S3 sync)
2. Adapt `rclone.conf` to your needs

```bash
docker-compose up --build -d
docker-compose exec -d rclone rclone --rc --rc-addr "0.0.0.0:5572" -v sync source:bucket target:bucket
```

3. Visit http://localhost:8080 in your browser for the exporter, or http://localhost:9090 for Prometheus

## CLI args

```console
rclone-exporter (version <version>, <commit>, <date>)

All flags can be read from Environment variables as well (replace . with _ , e.g. LOG_LEVEL).
However, CLI flags take precedence.

      --bindAddr string                    IP Address to bind to listen for Prometheus scrapes (default "0.0.0.0:8080")
      --log.level string                   Logging level (default "info")
      --push.interval string               Interval of push metrics from rclone. Accepts Go time formats (e.g. 30s). 0 disables regular pushes (default "0")
      --push.jobname string                Job name when pushed. Defaults to hostname. (default "<hostname>")
      --push.url string                    Pushgateway URL to push metrics (e.g. http://pushgateway:9091)
      --scrape.basicauth.password string   Password if rclone instance is BasicAuth protected (ENV var preferred)
      --scrape.basicauth.username string   Username if rclone instance is BasicAuth protected (ENV var preferred)
      --scrape.url string                  Base URL of the rclone instance with rc enabled (default "http://localhost:5572/")
```
