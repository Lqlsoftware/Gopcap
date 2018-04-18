package gopcap

import (
	"github.com/google/gopacket/pcap"
	"fmt"
)

func getAdapter() *pcap.Interface {
	adapters,err := pcap.FindAllDevs()
	check(err)
	idx := 0
	for idx = range adapters {
		if adapters[idx].Addresses != nil {
			break
		}
	}
	fmt.Println("Select Devices:",adapters[idx].Name)
	fmt.Println("Address:",adapters[idx].Addresses[len(adapters[idx].Addresses) - 1].IP)
	return &adapters[idx]
}