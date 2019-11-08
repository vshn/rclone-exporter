# rclone-exporter

A prometheus metrics exporter for Rclone

## Build

Auto-build is enabled on Docker Hub upon git tag pushes.

### With Go

```bash
go mod download
go build ./
```

### With Docker

```bash
tag=docker.io/vshn/rclone-exporter:latest
docker build -t docker.io/vshn/rclone-exporter:latest .
docker run --rm -p 8080:8080 ${tag} <CLI-args-see-below>
```

## Run

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

3. Visit `http://localhost:8080` in your browser for the exporter, or `http://localhost:9090` for Prometheus

## CLI args

```console
Usage of /usr/bin/rclone_exporter:
      --bindAddr string        IP Address to bind to listen for Prometheus scrapes (default "0.0.0.0:8080")
      --log.level string       Logging level (default "info")
      --push.interval string   Interval of push metrics from rclone. Accepts Go time formats (e.g. 30s). 0 disables regular pushes (default "0")
      --push.jobname string    Job name when pushed. Defaults to hostname. (default "e9b1544ef8d8")
      --push.url string        Pushgateway URL to push metrics (e.g. http://pushgateway:9091)
      --scrape.url string      Base URL of the rclone instance with rc enabled (default "http://localhost:5572/")
pflag: help requested
```
