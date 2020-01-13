package collector

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/vexxhost/network-exporter/network_api"
)

var bgpExpected = `
# HELP network_bgp_messages_received BGP messages received
# TYPE network_bgp_messages_received gauge
network_bgp_messages_received{as="65001",device="fake-dev",remote="10.10.10.1"} 1234
network_bgp_messages_received{as="65002",device="fake-dev",remote="10.10.10.2"} 0
# HELP network_bgp_messages_sent BGP messages sent
# TYPE network_bgp_messages_sent gauge
network_bgp_messages_sent{as="65001",device="fake-dev",remote="10.10.10.1"} 6789
network_bgp_messages_sent{as="65002",device="fake-dev",remote="10.10.10.2"} 0
# HELP network_bgp_prefix_count BGP prefix count
# TYPE network_bgp_prefix_count gauge
network_bgp_prefix_count{as="65001",device="fake-dev",remote="10.10.10.1"} 1337
network_bgp_prefix_count{as="65002",device="fake-dev",remote="10.10.10.2"} 0
# HELP network_bgp_up BGP session is established (up = 1)
# TYPE network_bgp_up gauge
network_bgp_up{as="65001",device="fake-dev",remote="10.10.10.1"} 1
network_bgp_up{as="65002",device="fake-dev",remote="10.10.10.2"} 0
`

type FakeBgpApi struct {
	network_api.API
}

func (a *FakeBgpApi) Peers() ([]network_api.BgpPeer, error) {
	return []network_api.BgpPeer{
		network_api.BgpPeer{
			Up:               true,
			AS:               65001,
			Remote:           "10.10.10.1",
			PrefixCount:      float64(1337),
			MessagesReceived: float64(1234),
			MessagesSent:     float64(6789),
		},
		network_api.BgpPeer{
			Up:               false,
			AS:               65002,
			Remote:           "10.10.10.2",
			PrefixCount:      float64(0),
			MessagesReceived: float64(0),
			MessagesSent:     float64(0),
		},
	}, nil
}

func TestBGPCollect(t *testing.T) {
	collector := NewBgpCollector("fake-dev", &FakeBgpApi{})

	err := testutil.CollectAndCompare(collector, strings.NewReader(bgpExpected))
	assert.NoError(t, err)
}
