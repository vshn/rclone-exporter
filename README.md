# rclone-exporter (UNMAINTAINED)

## This project is not maintained anymore

Rclone features prometheus metrics built-in, see https://rclone.org/flags/

---

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
tag=vshn/rclone-exporter:v1
docker build -t ${tag} .
docker run --rm -p 8080:8080 ${tag} <CLI-args-see-below>
```

## Run

(Hint: When running with BasicAuth, the logs will not output the credentials. 
This might be confusing when trying to debug an issue)

### Run just the exporter

```bash
docker run --rm -p 8080:8080 ${tag} <CLI-args-see-below>
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

## Environment variables

Alternatively, all flags are also configurable with Environment variables.
Replace the `.` char with `_` and uppercase the names in order for them to be recognized.

e.g. `--log.level debug` becomes `LOG_LEVEL=debug`
