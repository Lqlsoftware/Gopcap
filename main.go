package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"fmt"
	"strconv"
)

var port layers.TCPPort
var handle *pcap.Handle
var chMap = make(map[layers.TCPPort]chan gopacket.Packet)
var useMap = make(chan bool,1)

func main() {
	// 开启监听
	listener(8998)
}

// 监听TCP端口
func listener(portNo layers.TCPPort) {
	devices,err := pcap.FindAllDevs()
	check(err)
	idx := 0
	for idx = range devices {
		if devices[idx].Addresses != nil {
			break
		}
	}
	fmt.Println("Select Devices:",devices[idx].Name)
	fmt.Println("Address:",devices[idx].Addresses)
	fmt.Println("Port:",portNo)
	fmt.Print("Start listening\n\n")
	port = portNo
	handle,err = pcap.OpenLive(devices[idx].Name, 65535, true, pcap.BlockForever)
	check(err)
	defer handle.Close()
	err = handle.SetBPFFilter("tcp and dst port " + strconv.Itoa(int(portNo)))
	check(err)
	//Create a new PacketDataSource
	src := gopacket.NewPacketSource(handle, handle.LinkType())
	//Packets returns a channel of packets
	in := src.Packets()
	for true {
		select {
		case packet := <-in:
			// 交给处理函数
			packetHandler(packet)
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}