package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type (
	CoreStats struct {
		Speed       float64
		Bytes       float64
		Errors      float64
		FatalError  bool
		RetryError  bool
		Checks      float64
		Transfers   float64
		Deletes     float64
		ElapsedTime float64
	}
)

var (
	namespace   = "rclone"
	bytesMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "bytes_transferred_total",
		Help:      "total transferred bytes since the start of the rclone process",
	})
	speedMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "speed",
		Help:      "average speed in bytes/sec since start of the rclone process",
	})
	errorsMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "errors_total",
		Help:      "number of errors",
	})
	checksMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "checked_files_total",
		Help:      "number of checked files",
	})
	transfersMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "files_transferred_total",
		Help:      "number of transferred files",
	})
	deletesMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "files_deleted_total",
		Help:      "number of deleted files",
	})
	fatalErrorMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "fatal_error",
		Help:      "whether there has been at least one FatalError",
	})
	retryErrorMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "retry_error",
		Help:      "whether there has been at least one non-NoRetryError",
	})

	upMetrics = promauto.NewGauge(prometheus.GaugeOpts{
		Name:      "up",
		Help:      "whether the last scrape was successful",
	})
)
