package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/vexxhost/network-exporter/network_api"
)

type DeviceCollector struct {
	prometheus.Collector

	API network_api.API

	Up   *prometheus.Desc
	Info *prometheus.Desc

	BootTimestamp *prometheus.Desc

	FreeMemory  *prometheus.Desc
	TotalMemory *prometheus.Desc
}

func NewDeviceCollector(device string, api network_api.API) *DeviceCollector {
	return &DeviceCollector{
		API: api,

		Up: prometheus.NewDesc(
			"network_device_up",
			"Network device up",
			[]string{}, prometheus.Labels{"device": device},
		),
		Info: prometheus.NewDesc(
			"network_device_info",
			"Network device info",
			[]string{"model", "serial", "version"}, prometheus.Labels{"device": device},
		),

		BootTimestamp: prometheus.NewDesc(
			"network_device_boot_timestamp",
			"Network device boot timestamp",
			[]string{}, prometheus.Labels{"device": device},
		),

		FreeMemory: prometheus.NewDesc(
			"network_device_free_memory",
			"Network device free memory",
			[]string{}, prometheus.Labels{"device": device},
		),
		TotalMemory: prometheus.NewDesc(
			"network_device_total_memory",
			"Network device total memory",
			[]string{}, prometheus.Labels{"device": device},
		),
	}
}

func (c *DeviceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
	ch <- c.Info
	ch <- c.BootTimestamp
	ch <- c.FreeMemory
	ch <- c.TotalMemory
}

func (c *DeviceCollector) Collect(ch chan<- prometheus.Metric) {
	info, err := c.API.Info()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			c.Up, prometheus.GaugeValue,
			float64(0),
		)

		log.Errorln(err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.Up, prometheus.GaugeValue,
		float64(1),
	)

	ch <- prometheus.MustNewConstMetric(
		c.Info, prometheus.GaugeValue,
		float64(1),
		info.Model, info.Serial, info.Version,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BootTimestamp, prometheus.GaugeValue,
		info.BootTimestamp,
	)

	ch <- prometheus.MustNewConstMetric(
		c.FreeMemory, prometheus.GaugeValue,
		info.FreeMemory,
	)
	ch <- prometheus.MustNewConstMetric(
		c.TotalMemory, prometheus.GaugeValue,
		info.TotalMemory,
	)
}
