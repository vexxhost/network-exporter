package main

import (
	"net/http"
	"os"

	"github.com/vexxhost/network-exporter/collector"
	"github.com/vexxhost/network-exporter/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(":9615").String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		configFile = kingpin.Flag(
			"config.file",
			"Configuration file path",
		).Default("config.yaml").String()
	)

	kingpin.Version(version.Print("network-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	file, err := os.Open(*configFile)
	if err != nil {
		log.Fatalln(err)
	}

	config, err := config.Load(file)
	if err != nil {
		log.Fatalln(err)
	}

	for _, device := range config.Devices {
		api := device.API()

		bgp := collector.NewBgpCollector(device.Name, api)
		prometheus.MustRegister(bgp)

		deviceInfo := collector.NewDeviceCollector(device.Name, api)
		prometheus.MustRegister(deviceInfo)

		iface := collector.NewInterfaceCollector(device.Name, api)
		prometheus.MustRegister(iface)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Network Exporter</title></head>
			<body>
			<h1>Network Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
