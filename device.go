package gopcap

import (
	"github.com/google/gopacket/pcap"
	"log"
)

// 自动选择网络适配器
func getAdapter() *pcap.Interface {
	adapters,err := pcap.FindAllDevs()
	check(err)
	idx := 0
	for idx = range adapters {
		if adapters[idx].Addresses != nil && len(adapters[idx].Addresses) != 0 {
			break
		}
	}
	// 输出IPv4地址
	log.Print("IPv4: ", getIPV4(&adapters[idx]))
	return &adapters[idx]
}