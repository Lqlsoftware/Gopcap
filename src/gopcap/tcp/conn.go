package tcp

import (
	"encoding/binary"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Connection struct {
	// 本机信息
	srcIP		[]byte
	srcPort		layers.TCPPort
	srcMac		[]byte
	srcSeq		uint32
	// 客户端信息
	dstIP		[]byte
	dstPort		layers.TCPPort
	dstMac		[]byte
	dstSeq		uint32
	dstAck		uint32
	dstMSS		uint16
	Channel		*chan gopacket.Packet
	State		State
}


// 新建连接
func NewConnection(channel *chan gopacket.Packet, request gopacket.Packet) *Connection {
	reqTCP := request.Layer(layers.LayerTypeTCP).(*layers.TCP)
	reqIP := request.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	reqETH := request.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)
	MSS := uint16(1469)
	for _,v := range reqTCP.Options {
		if v.OptionType == layers.TCPOptionKindMSS {
			MSS = binary.BigEndian.Uint16(v.OptionData)
		}
	}
	// 新建连接
	conn := &Connection{
		srcIP:		reqIP.DstIP,
		srcPort:	reqTCP.DstPort,
		srcMac:		reqETH.DstMAC,
		srcSeq:		0,
		dstIP:		reqIP.SrcIP,
		dstPort:	reqTCP.SrcPort,
		dstMac:		reqETH.SrcMAC,
		dstSeq:		reqTCP.Seq,
		dstMSS:		MSS,
		Channel:	channel,
		State: 		UNCONNECT,
	}
	// 进行TCP握手
	// 发送SYN
	conn.sendSYN()
	conn.State = WAITSYNACK
	return conn
}

// 更新Connect数据
func (conn *Connection)Update(rawPacket gopacket.Packet) {
	tcp := rawPacket.Layer(layers.LayerTypeTCP).(*layers.TCP)
	conn.srcSeq = tcp.Ack
	conn.dstSeq = tcp.Seq + uint32(len(tcp.Payload))
	conn.dstAck = tcp.Ack
}