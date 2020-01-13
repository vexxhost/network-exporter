package collector

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/vexxhost/network-exporter/network_api"
)

var deviceExpected = `
# HELP network_device_free_memory Network device free memory
# TYPE network_device_free_memory gauge
network_device_free_memory{device="fake-dev"} 1.024e+06
# HELP network_device_info Network device info
# TYPE network_device_info gauge
network_device_info{device="fake-dev",model="fast",serial="1337",version="1.33.7"} 1
# HELP network_device_total_memory Network device total memory
# TYPE network_device_total_memory gauge
network_device_total_memory{device="fake-dev"} 2.048e+06
# HELP network_device_up Network device up
# TYPE network_device_up gauge
network_device_up{device="fake-dev"} 1
# HELP network_device_uptime Network device uptime
# TYPE network_device_uptime gauge
network_device_uptime{device="fake-dev"} 542594
`

type FakeInfoApi struct {
	network_api.API
}

func (a *FakeInfoApi) Info() (*network_api.Info, error) {
	return &network_api.Info{
		Uptime:      542594,
		FreeMemory:  1024000,
		TotalMemory: 2048000,
		Model:       "fast",
		Serial:      "1337",
		Version:     "1.33.7",
	}, nil
}

func TestDeviceCollect(t *testing.T) {
	collector := NewDeviceCollector("fake-dev", &FakeInfoApi{})

	err := testutil.CollectAndCompare(collector, strings.NewReader(deviceExpected))
	assert.NoError(t, err)
}
