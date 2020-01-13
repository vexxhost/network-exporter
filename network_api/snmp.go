package network_api

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/soniah/gosnmp"
)

type SnmpAPI struct {
	API

	Target    string
	Community string

	GoSnmp *gosnmp.GoSNMP
}

func NewSnmpAPI(target string, community string) *SnmpAPI {
	return &SnmpAPI{
		Target:    target,
		Community: community,
		GoSnmp: &gosnmp.GoSNMP{
			Target:             target,
			Port:               161,
			Transport:          "udp",
			Community:          community,
			Version:            gosnmp.Version2c,
			Timeout:            time.Duration(2) * time.Second,
			Retries:            3,
			ExponentialTimeout: true,
			MaxOids:            gosnmp.MaxOids,
		},
	}
}

func (a *SnmpAPI) getClient() *gosnmp.GoSNMP {
	return &gosnmp.GoSNMP{
		Target:             a.Target,
		Port:               161,
		Transport:          "udp",
		Community:          a.Community,
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(2) * time.Second,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}
}

func (a *SnmpAPI) Info() (*Info, error) {
	info := &Info{}

	// Uptime: 1.3.6.1.2.1.1.3
	walks, err := a.walk([]string{
		"1.3.6.1.2.1.1",
		"1.0.8802.1.1.2.1.3.4",
	})
	if err != nil {
		return nil, err
	}

	for _, pdus := range walks {
		for _, pdu := range pdus {
			fmt.Println(pdu.Name)

			// Uptime
			if pdu.Name == ".1.3.6.1.2.1.1.3.0" {
				info.Uptime = float64(pdu.Value.(uint))
			}

			// Model
			if pdu.Name == ".1.0.8802.1.1.2.1.3.4.0" {
				info.Model = string(pdu.Value.([]byte))
			}
		}
	}

	return info, nil
}

func (a *SnmpAPI) Peers() ([]BgpPeer, error) {
	return nil, errors.New("Unimplemented")
}

type SnmpWalkResult struct {
	OID   string
	PDUs  []gosnmp.SnmpPDU
	Error error
}

func (a *SnmpAPI) walk(oids []string) (map[string][]gosnmp.SnmpPDU, error) {
	var err error
	walks := map[string][]gosnmp.SnmpPDU{}

	pdus := make(chan SnmpWalkResult, len(oids))
	errors := make(chan error, len(oids))

	wg := sync.WaitGroup{}
	for _, oid := range oids {
		wg.Add(1)
		go func(oid string) {
			defer wg.Done()

			client := a.getClient()
			client.Connect()

			pduData, err := client.WalkAll(oid)
			errors <- err
			pdus <- SnmpWalkResult{
				OID:   oid,
				PDUs:  pduData,
				Error: err,
			}
		}(oid)
	}

	wg.Wait()

	close(pdus)
	close(errors)

	for r := range pdus {
		walks[r.OID] = r.PDUs
	}

	for e := range errors {
		if e != nil {
			return nil, err
		}
	}

	return walks, err
}

func (a *SnmpAPI) Interfaces() ([]Interface, error) {
	interfaces := []Interface{}

	walks, err := a.walk([]string{
		"1.3.6.1.2.1.31.1.1.1.1",
		"1.3.6.1.2.1.31.1.1.1.18",
		"1.3.6.1.2.1.2.2.1.3",
		"1.3.6.1.2.1.31.1.1.1.6",
		"1.3.6.1.2.1.2.2.1.11",
		"1.3.6.1.2.1.2.2.1.12",
		"1.3.6.1.2.1.2.2.1.13",
		"1.3.6.1.2.1.2.2.1.14",
		"1.3.6.1.2.1.31.1.1.1.10",
		"1.3.6.1.2.1.2.2.1.17",
		"1.3.6.1.2.1.2.2.1.18",
		"1.3.6.1.2.1.2.2.1.19",
		"1.3.6.1.2.1.2.2.1.20",
	})
	if err != nil {
		return nil, err
	}

	ifaceMap := map[string]*Interface{}

	// NOTE(mnaser): Create interfaces for L2 interfaces only
	for _, pdu := range walks["1.3.6.1.2.1.2.2.1.3"] {
		s := strings.Split(pdu.Name, ".")
		index := s[len(s)-1]

		if pdu.Value.(int) == 6 {
			ifaceMap[index] = &Interface{}
		}
	}

	for oid, pdus := range walks {
		for _, pdu := range pdus {
			s := strings.Split(pdu.Name, ".")
			index := s[len(s)-1]

			// Skip if we don't have a map (means non L2 device)
			if _, ok := ifaceMap[index]; !ok {
				continue
			}

			// Name
			if oid == "1.3.6.1.2.1.31.1.1.1.1" {
				ifaceMap[index].Name = string(pdu.Value.([]byte))
			}
			// Description
			if oid == "1.3.6.1.2.1.31.1.1.1.18" {
				ifaceMap[index].Description = string(pdu.Value.([]byte))
			}

			// RxPackets
			if oid == "1.3.6.1.2.1.2.2.1.11" || oid == "1.3.6.1.2.1.2.2.1.12" {
				ifaceMap[index].RxPackets += float64(pdu.Value.(uint))
			}
			// RxBits
			if oid == "1.3.6.1.2.1.31.1.1.1.6" {
				ifaceMap[index].RxBits = float64(pdu.Value.(uint64)) * 8
			}
			// RxDrops
			if oid == "1.3.6.1.2.1.2.2.1.13" {
				ifaceMap[index].RxDrops = float64(pdu.Value.(uint))
			}
			// RxErrors
			if oid == "1.3.6.1.2.1.2.2.1.14" {
				ifaceMap[index].RxErrors = float64(pdu.Value.(uint))
			}

			// TxPackets
			if oid == "1.3.6.1.2.1.2.2.1.17" || oid == "1.3.6.1.2.1.2.2.1.18" {
				ifaceMap[index].TxPackets += float64(pdu.Value.(uint))
			}
			// TxBits
			if oid == "1.3.6.1.2.1.31.1.1.1.10" {
				ifaceMap[index].TxBits = float64(pdu.Value.(uint64)) * 8
			}
			// TxDrops
			if oid == "1.3.6.1.2.1.2.2.1.19" {
				ifaceMap[index].TxDrops = float64(pdu.Value.(uint))
			}
			// TxErrors
			if oid == "1.3.6.1.2.1.2.2.1.20" {
				ifaceMap[index].TxErrors = float64(pdu.Value.(uint))
			}
		}
	}

	for _, iface := range ifaceMap {
		interfaces = append(interfaces, *iface)
	}

	return interfaces, nil
}
