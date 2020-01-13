package network_api

import (
	"time"

	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"
)

type AristaAPI struct {
	API

	Transport string
	Host      string
	Username  string
	Password  string
	Port      int
}

func NewAristaAPI(transport string, host string, username string, password string, port int) API {
	return &AristaAPI{
		Transport: transport,
		Host:      host,
		Username:  username,
		Password:  password,
		Port:      port,
	}
}

func (a *AristaAPI) getNode() (*goeapi.Node, error) {
	return goeapi.Connect(a.Transport, a.Host, a.Username, a.Password, a.Port)
}

func (a *AristaAPI) Info() (*Info, error) {
	node, err := a.getNode()
	if err != nil {
		return nil, err
	}

	var version module.ShowVersion
	handle, _ := node.GetHandle("json")
	if err := handle.Enable(&version); err != nil {
		return nil, err
	}

	return &Info{
		Uptime: float64(time.Now().Unix()) - version.BootupTimestamp,

		FreeMemory:  float64(version.MemFree),
		TotalMemory: float64(version.MemTotal),

		Model:   version.ModelName,
		Serial:  version.SerialNumber,
		Version: version.Version,
	}, nil
}

func (a *AristaAPI) Interfaces() ([]Interface, error) {
	node, err := a.getNode()
	if err != nil {
		return nil, err
	}

	interfaces := []Interface{}

	show := module.Show(node)
	ifaceData := show.ShowInterfaces()

	for _, iface := range ifaceData.Interfaces {
		counters := iface.InterfaceCounters

		rxPackets := float64(counters.InBroadcastPkts) + float64(counters.InMulticastPkts) + float64(counters.InUcastPkts)
		rxBits := float64(counters.InOctets) * 8
		rxDrops := float64(counters.InDiscards)
		rxErrors := float64(counters.TotalInErrors)

		txPackets := float64(counters.OutBroadcastPkts) + float64(counters.OutMulticastPkts) + float64(counters.OutUcastPkts)
		txBits := float64(counters.OutOctets) * 8
		txDrops := float64(counters.OutDiscards)
		txErrors := float64(counters.TotalOutErrors)

		iface := Interface{
			Name:        iface.Name,
			Description: iface.Description,

			RxPackets: rxPackets,
			RxBits:    rxBits,
			RxDrops:   rxDrops,
			RxErrors:  rxErrors,

			TxPackets: txPackets,
			TxBits:    txBits,
			TxDrops:   txDrops,
			TxErrors:  txErrors,
		}

		interfaces = append(interfaces, iface)
	}

	return interfaces, nil
}

func (a *AristaAPI) Peers() ([]BgpPeer, error) {
	node, err := a.getNode()
	if err != nil {
		return nil, err
	}

	peers := []BgpPeer{}

	show := module.Show(node)
	summary, err := show.ShowIPBGPSummary()
	if err != nil {
		return peers, err
	}

	for _, vrf := range summary.VRFs {
		for peerId, peer := range vrf.Peers {
			peer := BgpPeer{
				Up:               (peer.PeerState == "Established"),
				AS:               peer.ASN,
				Remote:           peerId,
				PrefixCount:      float64(peer.PrefixAccepted),
				MessagesReceived: float64(peer.MsgReceived),
				MessagesSent:     float64(peer.MsgSent),
			}

			peers = append(peers, peer)
		}
	}

	return peers, nil
}
