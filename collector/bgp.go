package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/vexxhost/network-exporter/network_api"
)

type BgpCollector struct {
	prometheus.Collector

	API network_api.API

	PrefixCount      *prometheus.Desc
	PeerUp           *prometheus.Desc
	MessagesReceived *prometheus.Desc
	MessagesSent     *prometheus.Desc
}

func NewBgpCollector(device string, api network_api.API) *BgpCollector {
	return &BgpCollector{
		API: api,

		PrefixCount: prometheus.NewDesc(
			"network_bgp_prefix_count",
			"BGP prefix count",
			[]string{"as", "remote"}, prometheus.Labels{"device": device},
		),
		PeerUp: prometheus.NewDesc(
			"network_bgp_up",
			"BGP session is established (up = 1)",
			[]string{"as", "remote"}, prometheus.Labels{"device": device},
		),
		MessagesReceived: prometheus.NewDesc(
			"network_bgp_messages_received",
			"BGP messages received",
			[]string{"as", "remote"}, prometheus.Labels{"device": device},
		),
		MessagesSent: prometheus.NewDesc(
			"network_bgp_messages_sent",
			"BGP messages sent",
			[]string{"as", "remote"}, prometheus.Labels{"device": device},
		),
	}
}

func (c *BgpCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.PrefixCount
	ch <- c.PeerUp
	ch <- c.MessagesReceived
	ch <- c.MessagesSent
}

func (c *BgpCollector) Collect(ch chan<- prometheus.Metric) {
	peers, err := c.API.Peers()
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, peer := range peers {
		up := 0
		if peer.Up {
			up = 1
		}

		ch <- prometheus.MustNewConstMetric(
			c.PrefixCount, prometheus.GaugeValue,
			peer.PrefixCount,
			strconv.FormatInt(peer.AS, 10), peer.Remote,
		)
		ch <- prometheus.MustNewConstMetric(
			c.PeerUp, prometheus.GaugeValue,
			float64(up),
			strconv.FormatInt(peer.AS, 10), peer.Remote,
		)
		ch <- prometheus.MustNewConstMetric(
			c.MessagesReceived, prometheus.GaugeValue,
			peer.MessagesReceived,
			strconv.FormatInt(peer.AS, 10), peer.Remote,
		)
		ch <- prometheus.MustNewConstMetric(
			c.MessagesSent, prometheus.GaugeValue,
			peer.MessagesSent,
			strconv.FormatInt(peer.AS, 10), peer.Remote,
		)
	}
}
