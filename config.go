package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	url2 "net/url"
	"os"
	"strings"
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
			BasicAuth: BasicAuthMap{
				Username: "",
				Password: "",
			},
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

	//viper.SetEnvPrefix("RE")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	defaults := CreateDefaultConfig()
	ser, err := json.Marshal(defaults)
	if err == nil {
		viper.SetConfigType("ser")
		return viper.ReadConfig(bytes.NewBuffer(ser))
	} else {
		return err
	}
}

func setupFlags() {
	cfg := CreateDefaultConfig()

	flag.String("push.url", cfg.Push.Url, "Pushgateway URL to push metrics (e.g. http://pushgateway:9091)")
	flag.String("push.interval", cfg.Push.Interval, "Interval of push metrics from rclone. Accepts Go time formats (e.g. 30s). 0 disables regular pushes")
	flag.String("push.jobname", cfg.Push.JobName, fmt.Sprintf("Job name when pushed. Defaults to hostname."))
	flag.String("scrape.url", cfg.Scrape.Url, "Base URL of the rclone instance with rc enabled")
	flag.String("scrape.basicauth.username", cfg.Scrape.BasicAuth.Username, "Username if rclone instance is BasicAuth protected (ENV var preferred)")
	flag.String("scrape.basicauth.password", cfg.Scrape.BasicAuth.Password, "Password if rclone instance is BasicAuth protected (ENV var preferred)")
	flag.String("bindAddr", cfg.BindAddr, "IP Address to bind to listen for Prometheus scrapes")
	flag.String("log.level", cfg.Log.Level, "Logging level")

	if err := viper.BindPFlags(flag.CommandLine); err != nil {
		log.Fatal(err)
	}
}

func GetConfig() ConfigMap {
	cfg := CreateDefaultConfig()
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func GetFriendlyUrlString(url string) string {
	// we don't want to log the credentials
	parsed, err := url2.Parse(url)
	if err != nil {
		return err.Error()
	}
	return GetFriendlyUrl(parsed)
}

func GetFriendlyUrl(url *url2.URL) string {
	// we don't want to log the credentials
	temp := url.User
	url.User = nil
	defer func() { url.User = temp }()
	return url.String()
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
		Url       string
		BasicAuth BasicAuthMap
	}
	BasicAuthMap struct {
		Username string
		Password string
	}
	PushMap struct {
		Url      string
		Interval string
		JobName  string
	}
)
