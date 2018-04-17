package tcp

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)


// 发送ACK
func (conn *Connection)sendAck() {
	// ethernet
	ethLayer := layers.Ethernet{
		SrcMAC:		conn.srcMac,
		DstMAC:		conn.dstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}
	// ip
	ipLayer := layers.IPv4{
		SrcIP:    	conn.srcIP,
		DstIP:    	conn.dstIP,
		TTL:	  	64,
		Protocol: 	layers.IPProtocolTCP,
		Version:  	4,
		Flags:	  	2,
	}
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: 	conn.srcPort,
		DstPort: 	conn.dstPort,
		ACK:	 	true,
		Ack:	 	conn.dstSeq,
		Seq:	 	conn.srcSeq,
		Window:  	0xFFFF,
	}
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	buf := gopacket.NewSerializeBuffer()
	conn.writeRaw(buf, &ethLayer, &ipLayer, &tcpLayer)
}

// 发送SYN
func (conn *Connection)sendSYN() {
	// ethernet
	ethLayer := layers.Ethernet{
		SrcMAC:		conn.srcMac,
		DstMAC:		conn.dstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}
	// ip
	ipLayer := layers.IPv4{
		SrcIP:    	conn.srcIP,
		DstIP:    	conn.dstIP,
		TTL:	  	64,
		Protocol: 	layers.IPProtocolTCP,
		Version:  	4,
		Flags:	  	2,
	}
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
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	buf := gopacket.NewSerializeBuffer()
	conn.writeRaw(buf, &ethLayer, &ipLayer, &tcpLayer)
}

// 发送FIN
func (conn *Connection)sendFin() {
	// ethernet
	ethLayer := layers.Ethernet{
		SrcMAC:		conn.srcMac,
		DstMAC:		conn.dstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}
	// ip
	ipLayer := layers.IPv4{
		SrcIP:    	conn.srcIP,
		DstIP:    	conn.dstIP,
		TTL:	  	64,
		Protocol: 	layers.IPProtocolTCP,
		Version:  	4,
		Flags:	  	2,
	}
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
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	buf := gopacket.NewSerializeBuffer()
	conn.writeRaw(buf, &ethLayer, &ipLayer, &tcpLayer)
}