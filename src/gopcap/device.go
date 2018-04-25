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
		if adapters[idx].Addresses != nil {
			break
		}
	}
	log.Print("Devices: ",adapters[idx].Name)
	// 输出IPv4地址
	for _,v := range adapters[idx].Addresses {
		if len(v.IP) == 4 {
			log.Print("IPv4:    ",v.IP)
		}
	}
	return &adapters[idx]
}