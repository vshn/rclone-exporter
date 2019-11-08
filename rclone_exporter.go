package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	log "github.com/sirupsen/logrus"
	"net/http"
	url2 "net/url"
	"time"
)

func main() {

	if err := LoadConfig(); err != nil {
		log.WithError(err).Error("Could not load config.")
	}
	ConfigureLogging()

	cfg := GetConfig()

	url, err := url2.Parse(cfg.Scrape.Url)
	if err != nil {
		log.WithError(err).Fatal("Can not start with incorrect URL.")
	}
	url.Path = "/core/stats"
	coreStatsUrl = url.String()

	if cfg.Push.Url == "" {
		log.Warn("There is no pushgateway URL defined. Pushes will fail.")
	}

	log.WithFields(log.Fields{
		"scrape-url": cfg.Scrape.Url,
		"bind-addr":  cfg.BindAddr,
		"push-url":   cfg.Push.Url,
	}).Info("Starting exporter.")

	pusher = push.New(cfg.Push.Url, cfg.Push.JobName).Gatherer(prometheus.DefaultGatherer)

	if cfg.Push.Interval != "0" {
		go PushRegularly(cfg)
	}

	http.HandleFunc("/", handleMetrics)
	http.HandleFunc("/metrics", handleMetrics)
	http.HandleFunc("/push", handlePush)
	if err := http.ListenAndServe(cfg.BindAddr, nil); err != nil {
		log.Fatal(err)
	}
}

var (
	coreStatsUrl              = ""
	promHandler               = promhttp.Handler()
	pusher       *push.Pusher = nil
)

func PushRegularly(cfg ConfigMap) {
	if cfg.Push.Url == "" {
		log.Error("There is no pushgateway URL defined. Not pushing regularly.")
		return
	}
	interval, err := time.ParseDuration(cfg.Push.Interval)
	if err != nil {
		interval = time.Second * 30
		log.WithError(err).Warning("Pushing in default interval.")
	}
	log.WithField("interval", interval).Info("Pushing regularly to pushgateway.")
	for {
		scrapeErr := scrape()
		if scrapeErr != nil {
			log.WithError(err).Debug("Could not scrape metrics from rclone, assuming rclone is not running now.")
		}
		logEvent := log.WithField("url", cfg.Push.Url)
		err := pushMetrics(cfg)
		if err == nil {
			logEvent.Debug("Pushed metrics.")
		} else {
			logEvent.WithError(err).Error("Could not push to pushgateway.")
		}
		log.WithField("interval", interval).Debug("Wait for next push interval.")
		time.Sleep(interval)
	}
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	if err := scrape(); err != nil {
		log.Error(err)
	}
	promHandler.ServeHTTP(w, r)
}

func handlePush(w http.ResponseWriter, r *http.Request) {
	cfg := GetConfig()
	logEvent := log.WithField("url", cfg.Push.Url)
	err := pushMetrics(cfg)
	if err == nil {
		logEvent.Info("Pushed to pushgateway.")
		fmt.Fprintf(w, "Successfully pushed to %s", cfg.Push.Url)
	} else {
		logEvent.WithError(err).Error("Could not push to pushgateway.")
		w.WriteHeader(500)
		fmt.Fprintf(w, "Could not push to pushgateway: %s", err.Error())
	}
}

func scrape() error {
	coreStats := CoreStats{}
	log.WithField("url", coreStatsUrl).Debug("Collecting core stats.")
	if err := collect(coreStatsUrl, &coreStats); err != nil {
		resetStats()
		upMetrics.Set(0)
		return err
	} else {
		parseStats(coreStats)
		upMetrics.Set(1)
		return nil
	}
}

func pushMetrics(cfg ConfigMap) error {
	if cfg.Push.Url == "" {
		return errors.New("no pushgateway URL defined")
	}
	return pusher.Push()
}

func parseStats(s CoreStats) {
	speedMetric.Set(s.Speed)
	bytesMetric.Set(s.Bytes)
	errorsMetric.Set(s.Errors)
	retryErrorMetric.Set(0)
	fatalErrorMetric.Set(0)
	if s.FatalError {
		fatalErrorMetric.Set(1)
	}
	if s.RetryError {
		retryErrorMetric.Set(1)
	}
	checksMetric.Set(s.Checks)
	transfersMetric.Set(s.Transfers)
	deletesMetric.Set(s.Deletes)
}

func resetStats() {
	speedMetric.Set(0)
	bytesMetric.Set(0)
	errorsMetric.Set(0)
	retryErrorMetric.Set(0)
	fatalErrorMetric.Set(0)
	checksMetric.Set(0)
	transfersMetric.Set(0)
	deletesMetric.Set(0)
}

func collect(url string, model interface{}) error {
	response, err := http.PostForm(url, url2.Values{})
	if err == nil {
		decodeErr := json.NewDecoder(response.Body).Decode(model)
		return decodeErr
	} else {
		return err
	}
}
