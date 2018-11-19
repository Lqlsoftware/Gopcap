package tcp

import "github.com/google/gopacket/layers"

// 生成以太网/IP层的报头
func (connection *Connection)getUpperLayers() (*layers.Ethernet, *layers.IPv4) {
	// ethernet
	ethLayer := layers.Ethernet{
		SrcMAC:       connection.srcMac,
		DstMAC:       connection.dstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}
	// ip
	ipLayer := layers.IPv4{
		SrcIP:    connection.srcIP,
		DstIP:    connection.dstIP,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		Version:  4,
		Flags:    2,
	}
	return &ethLayer, &ipLayer
}

// 生成以太网报头
func (connection *Connection)getLayers() *layers.TCP {
	return &layers.TCP{
		SrcPort: connection.srcPort,
		DstPort: connection.dstPort,
		ACK:     true,
		Ack:     connection.dstSeq,
		Seq:     connection.srcSeq,
		Window:  0xFFFF,
	}
}