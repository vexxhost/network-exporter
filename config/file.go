package config

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/vexxhost/network-exporter/network_api"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Devices []DeviceConfig `yaml:"devices"`
}

type DeviceConfig struct {
	Type      string `yaml:"type"`
	Name      string `yaml:"name"`
	Transport string `yaml:"transport"`
	Hostname  string `yaml:"hostname"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

func Load(handle io.Reader) (*Config, error) {
	config := &Config{}

	data, err := ioutil.ReadAll(handle)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *DeviceConfig) String() string {
	return fmt.Sprintf(
		"[%s %s://%s@%s:%d]",
		c.Name,
		c.Transport, c.Username, c.Hostname, c.Port,
	)
}

func (c *DeviceConfig) API() network_api.API {

	switch c.Type {
	case "mikrotik":
		return network_api.NewMikrotikAPI(
			c.Hostname,
			c.Port,
			c.Username,
			c.Password,
		)
	case "arista":
		return network_api.NewAristaAPI(
			c.Transport,
			c.Hostname,
			c.Username,
			c.Password,
			c.Port,
		)
	}

	return nil
}
