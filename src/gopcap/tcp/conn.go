package tcp

import (
	"encoding/binary"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TCP Connection
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
	dstWin		uint16
	// 数据包channel
	Channel		*chan gopacket.Packet
	// 连接状态
	State		State
}


// 新建TCP连接
func NewConnection(channel *chan gopacket.Packet, request gopacket.Packet) *Connection {
	// 握手包TCP层
	reqTCP := request.Layer(layers.LayerTypeTCP).(*layers.TCP)
	// 握手包IP层
	reqIP := request.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	// 握手包以太网层
	reqETH := request.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)

	// 最大TCP段长度
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
		dstWin:		reqTCP.Window,
		Channel:	channel,
		State: 		UNCONNECT,
	}
	// 进行TCP握手 发送SYN
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
	conn.dstWin = tcp.Window
}