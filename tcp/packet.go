package tcp

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)


// 发送ACK
func (connection *Connection)sendAck() {
	ethLayer, ipLayer := connection.getUpperLayers()
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: connection.srcPort,
		DstPort: connection.dstPort,
		ACK:     true,
		Ack:     connection.dstSeq,
		Seq:     connection.srcSeq,
		Window:  0xFFFF,
	}
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	buf := gopacket.NewSerializeBuffer()
	connection.writeRaw(buf, ethLayer, ipLayer, &tcpLayer)
}

// 发送SYN
func (connection *Connection)sendSYN() {
	ethLayer, ipLayer := connection.getUpperLayers()
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: connection.srcPort,
		DstPort: connection.dstPort,
		SYN:     true,
		ACK:     true,
		Ack:     connection.dstSeq + 1,
		Window:  0xFFFF,
		Options: []layers.TCPOption{{layers.TCPOptionKindMSS,4,[]byte{5,189}}},
	}
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	buf := gopacket.NewSerializeBuffer()
	connection.writeRaw(buf, ethLayer, ipLayer, &tcpLayer)
}

// 发送FIN
func (connection *Connection)sendFin() {
	ethLayer, ipLayer := connection.getUpperLayers()
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: connection.srcPort,
		DstPort: connection.dstPort,
		ACK:     true,
		FIN:     true,
		Ack:     connection.dstSeq,
		Seq:     connection.srcSeq,
		Window:  0xFFFF,
	}
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	buf := gopacket.NewSerializeBuffer()
	connection.writeRaw(buf, ethLayer, ipLayer, &tcpLayer)
}