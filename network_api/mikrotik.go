package network_api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-routeros/routeros"
	"github.com/karrick/tparse"
)

type MikrotikAPI struct {
	API

	Address  string
	Username string
	Password string
}

func NewMikrotikAPI(hostname string, port int, username string, password string) API {
	address := fmt.Sprintf("%s:%d", hostname, port)

	return &MikrotikAPI{
		Address:  address,
		Username: username,
		Password: password,
	}
}

func (a *MikrotikAPI) getClient() (*routeros.Client, error) {
	return routeros.Dial(a.Address, a.Username, a.Password)
}

func (a *MikrotikAPI) Info() (*Info, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
	}

	systemData, err := client.Run(
		"/system/resource/print",
		"=.proplist=uptime,free-memory,total-memory,board-name,version",
	)
	if err != nil {
		return nil, err
	}

	boardData, err := client.Run(
		"/system/routerboard/print",
		"=.proplist=serial-number",
	)
	if err != nil {
		return nil, err
	}

	bootTimestamp, _ := tparse.ParseNow(time.RFC3339, "now-"+systemData.Re[0].Map["uptime"])

	freeMemory, _ := strconv.ParseFloat(systemData.Re[0].Map["free-memory"], 64)
	totalMemory, _ := strconv.ParseFloat(systemData.Re[0].Map["total-memory"], 64)

	return &Info{
		Uptime: float64(time.Now().Unix()) - float64(bootTimestamp.Unix()),

		FreeMemory:  freeMemory,
		TotalMemory: totalMemory,

		Model:   systemData.Re[0].Map["board-name"],
		Serial:  boardData.Re[0].Map["serial-number"],
		Version: systemData.Re[0].Map["version"],
	}, nil
}

func (a *MikrotikAPI) Interfaces() ([]Interface, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
	}

	interfaces := []Interface{}

	interfaceData, err := client.Run(
		"/interface/ethernet/print",
		"=.proplist=name,rx-packet,rx-byte,rx-drop,rx-error,tx-packet,tx-byte,tx-drop,tx-error",
	)
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaceData.Re {
		rxPackets, _ := strconv.ParseFloat(iface.Map["rx-packet"], 64)
		rxBytes, _ := strconv.ParseFloat(iface.Map["rx-byte"], 64)
		rxDrops, _ := strconv.ParseFloat(iface.Map["rx-drop"], 64)
		rxErrors, _ := strconv.ParseFloat(iface.Map["rx-error"], 64)

		txPackets, _ := strconv.ParseFloat(iface.Map["tx-packet"], 64)
		txBytes, _ := strconv.ParseFloat(iface.Map["tx-byte"], 64)
		txDrops, _ := strconv.ParseFloat(iface.Map["tx-drop"], 64)
		txErrors, _ := strconv.ParseFloat(iface.Map["tx-error"], 64)

		iface := Interface{
			Name:        iface.Map["name"],
			Description: iface.Map["comment"],

			RxPackets: rxPackets,
			RxBits:    rxBytes * 8,
			RxDrops:   rxDrops,
			RxErrors:  rxErrors,

			TxPackets: txPackets,
			TxBits:    txBytes * 8,
			TxDrops:   txDrops,
			TxErrors:  txErrors,
		}

		interfaces = append(interfaces, iface)
	}

	return interfaces, nil
}

func (a *MikrotikAPI) Peers() ([]BgpPeer, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
	}

	peers := []BgpPeer{}

	peerData, err := client.Run(
		"/routing/bgp/peer/print",
		"=.proplist=state,remote-as,remote-address,prefix-count,updates-sent,updates-received,withdrawn-sent,withdrawn-received",
	)
	if err != nil {
		return nil, err
	}

	for _, peer := range peerData.Re {
		asn, _ := strconv.ParseInt(peer.Map["remote-as"], 10, 64)
		prefixCount, _ := strconv.ParseFloat(peer.Map["prefix-count"], 64)

		updatesSent, _ := strconv.ParseFloat(peer.Map["updates-sent"], 64)
		updatesReceived, _ := strconv.ParseFloat(peer.Map["updates-received"], 64)
		withdrawnSent, _ := strconv.ParseFloat(peer.Map["withdrawn-sent"], 64)
		withdrawnReceived, _ := strconv.ParseFloat(peer.Map["withdrawn-received"], 64)

		messagesSent := updatesSent + withdrawnSent
		messagesReceived := updatesReceived + withdrawnReceived

		peer := BgpPeer{
			Up:               (peer.Map["state"] == "established"),
			AS:               asn,
			Remote:           peer.Map["remote-address"],
			PrefixCount:      prefixCount,
			MessagesReceived: messagesSent,
			MessagesSent:     messagesReceived,
		}

		peers = append(peers, peer)
	}

	return peers, nil
}
