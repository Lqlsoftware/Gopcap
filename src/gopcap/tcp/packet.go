package tcp

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)


// 发送ACK
func (conn *Connection)sendAck() {
	ethLayer, ipLayer := conn.getUpperLayers()
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: 	conn.srcPort,
		DstPort: 	conn.dstPort,
		ACK:	 	true,
		Ack:	 	conn.dstSeq,
		Seq:	 	conn.srcSeq,
		Window:  	0xFFFF,
	}
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	buf := gopacket.NewSerializeBuffer()
	conn.writeRaw(buf, ethLayer, ipLayer, &tcpLayer)
}

// 发送SYN
func (conn *Connection)sendSYN() {
	ethLayer, ipLayer := conn.getUpperLayers()
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: 	conn.srcPort,
		DstPort: 	conn.dstPort,
		SYN:     	true,
		ACK:	 	true,
		Ack:	 	conn.dstSeq + 1,
		Window:  	0xFFFF,
		Options:	[]layers.TCPOption{{layers.TCPOptionKindMSS,4,[]byte{5,189}}},
	}
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	buf := gopacket.NewSerializeBuffer()
	conn.writeRaw(buf, ethLayer, ipLayer, &tcpLayer)
}

// 发送FIN
func (conn *Connection)sendFin() {
	ethLayer, ipLayer := conn.getUpperLayers()
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: 	conn.srcPort,
		DstPort: 	conn.dstPort,
		ACK:	 	true,
		FIN:		true,
		Ack:	 	conn.dstSeq,
		Seq:	 	conn.srcSeq,
		Window:  	0xFFFF,
	}
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	buf := gopacket.NewSerializeBuffer()
	conn.writeRaw(buf, ethLayer, ipLayer, &tcpLayer)
}

// 发送默认ACK
func sendAck(packet gopacket.Packet) {
	// TCP层
	reqTCP := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
	// IP层
	reqIP := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	// 以太网层
	reqETH := packet.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)

	// 生成ACK 以太网层
	ethLayer := layers.Ethernet{
		SrcMAC:       reqETH.DstMAC,
		DstMAC:       reqETH.SrcMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}

	// 生成ACK IP层
	ipLayer := layers.IPv4{
		SrcIP:    reqIP.DstIP,
		DstIP:    reqIP.SrcIP,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		Version:  4,
		Flags:    2,
	}

	// 生成ACK TCP层
	tcpLayer := layers.TCP{
		SrcPort: 	reqTCP.DstPort,
		DstPort: 	reqTCP.SrcPort,
		ACK:	 	true,
		Ack:	 	reqTCP.Seq + 1,
		Seq:	 	reqTCP.Ack,
		Window:  	0xFFFF,
	}
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	buf := gopacket.NewSerializeBuffer()
	err := tcpLayer.SerializeTo(buf, gopacket.SerializeOptions{true,true})
	check(err)
	err = ipLayer.SerializeTo(buf, gopacket.SerializeOptions{true,true})
	check(err)
	err = ethLayer.SerializeTo(buf, gopacket.SerializeOptions{false,true})
	check(err)
	err = sendChannel.WritePacketData(buf.Bytes())
	check(err)
}