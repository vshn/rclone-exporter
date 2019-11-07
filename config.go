package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"os"
)

var defaultJobName = ""

func init() {
	defaultJobName, _ = os.Hostname()
	setupFlags()
}

func CreateDefaultConfig() ConfigMap {
	return ConfigMap{
		Log: LogMap{
			Level: "info",
		},
		Scrape: ScrapeMap{
			Url: "http://localhost:5572/",
		},
		Push: PushMap{
			Url:      "",
			Interval: "0",
			JobName:  defaultJobName,
		},
		BindAddr: "0.0.0.0:8080",
	}
}

func LoadConfig() error {
	flag.Parse()
	return nil
}

func setupFlags() {
	cfg := CreateDefaultConfig()
	flag.String("push.url", cfg.Push.Url, "Pushgateway URL to push metrics (e.g. http://pushgateway:9091)")
	flag.String("push.interval", cfg.Push.Interval, "Interval of push metrics from rclone. Accepts Go time formats (e.g. 30s). 0 disables regular pushes")
	flag.String("push.jobname", cfg.Push.JobName, fmt.Sprintf("Job name when pushed. Defaults to hostname."))
	flag.String("scrape.url", cfg.Scrape.Url, "Base URL of the rclone instance with rc enabled")
	flag.String("bindAddr", cfg.BindAddr, "IP Address to bind to listen for Prometheus scrapes")
	flag.String("log.level", cfg.Log.Level, "Logging level")
}

func GetConfig() ConfigMap {
	cfg := CreateDefaultConfig()
	cfg.Scrape.Url = flag.Lookup("scrape.url").Value.String()
	cfg.Push.Interval = flag.Lookup("push.interval").Value.String()
	cfg.Push.JobName = flag.Lookup("push.jobname").Value.String()
	cfg.Push.Url = flag.Lookup("push.url").Value.String()
	cfg.Log.Level = flag.Lookup("log.level").Value.String()
	cfg.BindAddr = flag.Lookup("bindAddr").Value.String()
	return cfg
}

func ConfigureLogging() {
	cfg := GetConfig()

	log.SetOutput(os.Stdout)

	level, err := log.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.WithField("error", err).Warn("Using info level.")
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
}

type (
	ConfigMap struct {
		Log      LogMap
		Scrape   ScrapeMap
		Push     PushMap
		BindAddr string
	}
	LogMap struct {
		Level     string
		Formatter string
	}
	ScrapeMap struct {
		Url string
	}
	PushMap struct {
		Url      string
		Interval string
		JobName  string
	}
)
