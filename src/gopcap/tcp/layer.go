package tcp

import "github.com/google/gopacket/layers"

// 生成以太网/IP层的报头
func (conn *Connection)getUpperLayers() (*layers.Ethernet, *layers.IPv4) {
	// ethernet
	ethLayer := layers.Ethernet{
		SrcMAC:       conn.srcMac,
		DstMAC:       conn.dstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}
	// ip
	ipLayer := layers.IPv4{
		SrcIP:    conn.srcIP,
		DstIP:    conn.dstIP,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		Version:  4,
		Flags:    2,
	}
	return &ethLayer, &ipLayer
}

// 生成以太网报头
func (conn *Connection)getLayers() *layers.TCP {
	return &layers.TCP{
		SrcPort: 	conn.srcPort,
		DstPort: 	conn.dstPort,
		ACK:	 	true,
		Ack:	 	conn.dstSeq,
		Seq:	 	conn.srcSeq,
		Window:  	0xFFFF,
	}
}