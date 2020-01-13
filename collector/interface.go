package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/vexxhost/network-exporter/network_api"
)

type InterfaceCollector struct {
	prometheus.Collector

	API network_api.API

	Info *prometheus.Desc

	RxPackets *prometheus.Desc
	RxBits    *prometheus.Desc
	RxDrops   *prometheus.Desc
	RxErrors  *prometheus.Desc

	TxPackets *prometheus.Desc
	TxBits    *prometheus.Desc
	TxDrops   *prometheus.Desc
	TxErrors  *prometheus.Desc
}

func NewInterfaceCollector(device string, api network_api.API) *InterfaceCollector {
	return &InterfaceCollector{
		API: api,

		Info: prometheus.NewDesc(
			"network_interface_info",
			"Interface information",
			[]string{"interface", "description"}, prometheus.Labels{"device": device},
		),

		RxPackets: prometheus.NewDesc(
			"network_interface_rx_packets",
			"Number of packets received",
			[]string{"interface"}, prometheus.Labels{"device": device},
		),
		RxBits: prometheus.NewDesc(
			"network_interface_rx_bits",
			"Number of bits received",
			[]string{"interface"}, prometheus.Labels{"device": device},
		),
		RxDrops: prometheus.NewDesc(
			"network_interface_rx_drops",
			"Number of receive packet drops",
			[]string{"interface"}, prometheus.Labels{"device": device},
		),
		RxErrors: prometheus.NewDesc(
			"network_interface_rx_errors",
			"Number of receive packet errors",
			[]string{"interface"}, prometheus.Labels{"device": device},
		),

		TxPackets: prometheus.NewDesc(
			"network_interface_tx_packets",
			"Number of packets transmitted",
			[]string{"interface"}, prometheus.Labels{"device": device},
		),
		TxBits: prometheus.NewDesc(
			"network_interface_tx_bits",
			"Number of bits transmitted",
			[]string{"interface"}, prometheus.Labels{"device": device},
		),
		TxDrops: prometheus.NewDesc(
			"network_interface_tx_drops",
			"Number of transmit packet drops",
			[]string{"interface"}, prometheus.Labels{"device": device},
		),
		TxErrors: prometheus.NewDesc(
			"network_interface_tx_errors",
			"Number of transmit packet errors",
			[]string{"interface"}, prometheus.Labels{"device": device},
		),
	}
}

func (c *InterfaceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Info

	ch <- c.RxPackets
	ch <- c.RxBits
	ch <- c.RxDrops
	ch <- c.RxErrors

	ch <- c.TxPackets
	ch <- c.TxBits
	ch <- c.TxDrops
	ch <- c.TxErrors
}

func (c *InterfaceCollector) Collect(ch chan<- prometheus.Metric) {
	interfaces, err := c.API.Interfaces()
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, iface := range interfaces {
		ch <- prometheus.MustNewConstMetric(
			c.Info, prometheus.GaugeValue,
			float64(1),
			iface.Name, iface.Description,
		)

		ch <- prometheus.MustNewConstMetric(
			c.RxPackets, prometheus.CounterValue,
			iface.RxPackets,
			iface.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.RxBits, prometheus.CounterValue,
			iface.RxBits,
			iface.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.RxDrops, prometheus.CounterValue,
			iface.RxDrops,
			iface.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.RxErrors, prometheus.CounterValue,
			iface.RxErrors,
			iface.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.TxPackets, prometheus.CounterValue,
			iface.TxPackets,
			iface.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TxBits, prometheus.CounterValue,
			iface.TxBits,
			iface.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TxDrops, prometheus.CounterValue,
			iface.TxDrops,
			iface.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TxErrors, prometheus.CounterValue,
			iface.TxErrors,
			iface.Name,
		)
	}
}
