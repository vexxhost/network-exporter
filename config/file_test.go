package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigWithNoDevices(t *testing.T) {
	configFile := `
---
devices: []
`

	config, err := Load(strings.NewReader(configFile))
	assert.NoError(t, err)

	assert.Empty(t, config.Devices)
}

func TestLoadConfig(t *testing.T) {
	configFile := `
---
devices:
- type: mikrotik
  name: switch-1
  community: prometheus
  transport: https
  hostname: switch.ip
  port: 443
  username: foo
  password: bar
`

	config, err := Load(strings.NewReader(configFile))
	assert.NoError(t, err)

	if assert.Len(t, config.Devices, 1) {
		assert.Equal(t, "mikrotik", config.Devices[0].Type)
		assert.Equal(t, "switch-1", config.Devices[0].Name)
		assert.Equal(t, "prometheus", config.Devices[0].Community)
		assert.Equal(t, "https", config.Devices[0].Transport)
		assert.Equal(t, "switch.ip", config.Devices[0].Hostname)
		assert.Equal(t, 443, config.Devices[0].Port)
		assert.Equal(t, "foo", config.Devices[0].Username)
		assert.Equal(t, "bar", config.Devices[0].Password)

		assert.Equal(t, "[switch-1 https://foo@switch.ip:443]", config.Devices[0].String())

	}
}
