package collector

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/vexxhost/network-exporter/network_api"
)

var interfaceExpected = `
# HELP network_interface_info Interface information
# TYPE network_interface_info gauge
network_interface_info{description="bad server",device="fake-dev",interface="Ethernet2"} 1
network_interface_info{description="uplink",device="fake-dev",interface="Ethernet1"} 1
# HELP network_interface_rx_bits Number of bits received
# TYPE network_interface_rx_bits counter
network_interface_rx_bits{device="fake-dev",interface="Ethernet1"} 54325
network_interface_rx_bits{device="fake-dev",interface="Ethernet2"} 1234
# HELP network_interface_rx_drops Number of receive packet drops
# TYPE network_interface_rx_drops counter
network_interface_rx_drops{device="fake-dev",interface="Ethernet1"} 0
network_interface_rx_drops{device="fake-dev",interface="Ethernet2"} 1
# HELP network_interface_rx_errors Number of receive packet errors
# TYPE network_interface_rx_errors counter
network_interface_rx_errors{device="fake-dev",interface="Ethernet1"} 0
network_interface_rx_errors{device="fake-dev",interface="Ethernet2"} 9
# HELP network_interface_rx_packets Number of packets received
# TYPE network_interface_rx_packets counter
network_interface_rx_packets{device="fake-dev",interface="Ethernet1"} 54325
network_interface_rx_packets{device="fake-dev",interface="Ethernet2"} 1234
# HELP network_interface_tx_bits Number of bits transmitted
# TYPE network_interface_tx_bits counter
network_interface_tx_bits{device="fake-dev",interface="Ethernet1"} 7443533453
network_interface_tx_bits{device="fake-dev",interface="Ethernet2"} 6345
# HELP network_interface_tx_drops Number of transmit packet drops
# TYPE network_interface_tx_drops counter
network_interface_tx_drops{device="fake-dev",interface="Ethernet1"} 0
network_interface_tx_drops{device="fake-dev",interface="Ethernet2"} 4
# HELP network_interface_tx_errors Number of transmit packet errors
# TYPE network_interface_tx_errors counter
network_interface_tx_errors{device="fake-dev",interface="Ethernet1"} 0
network_interface_tx_errors{device="fake-dev",interface="Ethernet2"} 6
# HELP network_interface_tx_packets Number of packets transmitted
# TYPE network_interface_tx_packets counter
network_interface_tx_packets{device="fake-dev",interface="Ethernet1"} 3454534
network_interface_tx_packets{device="fake-dev",interface="Ethernet2"} 9876
`

type FakeInterfaceApi struct {
	network_api.API
}

func (a *FakeInterfaceApi) Interfaces() ([]network_api.Interface, error) {
	return []network_api.Interface{
		network_api.Interface{
			Name:        "Ethernet1",
			Description: "uplink",
			RxPackets:   54325,
			RxBits:      6452523,
			RxDrops:     0,
			RxErrors:    0,
			TxPackets:   3454534,
			TxBits:      7443533453,
			TxDrops:     0,
			TxErrors:    0,
		},
		network_api.Interface{
			Name:        "Ethernet2",
			Description: "bad server",
			RxPackets:   1234,
			RxBits:      2345,
			RxDrops:     1,
			RxErrors:    9,
			TxPackets:   9876,
			TxBits:      6345,
			TxDrops:     4,
			TxErrors:    6,
		},
	}, nil
}

func TestInterfaceCollect(t *testing.T) {
	collector := NewInterfaceCollector("fake-dev", &FakeInterfaceApi{})

	err := testutil.CollectAndCompare(collector, strings.NewReader(interfaceExpected))
	assert.NoError(t, err)
}
