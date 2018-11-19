package tcp

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)


type Handler struct {
	In		*Receiver
	Out 	*Sender
}

func (handler Handler)Init(pipe *pcap.Handle) {
	handler.In.Init()
	handler.Out.Init(pipe)
}

func (handler Handler)Listen(In chan gopacket.Packet) {
	// listen all package
	for true {
		select {
		case packet := <-In:
			// 解析TCP报文
			TCPPacket := packet.TransportLayer().(*layers.TCP)


			// fetch from PortMap
			handler.In.PortMapLocker.RLock()
			pipe, exist := handler.In.PortMap[TCPPacket.SrcPort]
			handler.In.PortMapLocker.RUnlock()

			// new connection
			if !exist {
				IPPacket := packet.NetworkLayer().(*layers.IPv4)
				ETHPacket := packet.LinkLayer().(*layers.Ethernet)
				info := ConnectionInfo{ETHPacket.SrcMAC, ETHPacket.DstMAC, IPPacket.SrcIP, IPPacket.DstIP, TCPPacket.SrcPort, TCPPacket.SrcPort}
				pipe = *handler.In.AddPort(TCPPacket.SrcPort)
				go ProcessThread(&info, &pipe, &handler)
			}

			// push to port pipe
			pipe<- data
		}
	}
}

type ConnectionInfo struct {
	SrcMac		net.HardwareAddr
	DstMac		net.HardwareAddr
	SrcIP		net.IP
	DstIP		net.IP
	SrcPort		layers.TCPPort
	DstPort		layers.TCPPort
}