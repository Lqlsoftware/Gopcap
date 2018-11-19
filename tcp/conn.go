package tcp

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TCP State
type State uint8
const (
	CLOSED    	State = 0
	LISTEN		State = 1
	SYN_RECV	State = 2
	ESTABLISHED	State = 3
	CLOSE_WAIT	State = 4
	LAST_ACK	State = 5
)

const TIMEOUT = time.Millisecond * 500

func ProcessThread(info *ConnectionInfo, pipe *chan *layers.TCP, handler *Handler) {
	// defer remove port in map
	defer handler.In.RemovePort((*info).SrcPort)
	// new connection
	connection := NewConnection()
	for {
		select {
		case packet := <-(*pipe):
			connection.Receive(packet)
		}
	}
}

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
	// 连接状态
	State		State
}


// 新建连接
func NewConnection() *Connection {
	conn := &Connection{
		State:LISTEN,
	}
	return conn
}

func (connection *Connection)Receive(packet *layers.TCP) {
	switch connection.State {
	case CLOSED:
		// Reset
		connection.RST()
		connection.State = LISTEN
	case LISTEN:
		// SYN -> shakehands
		if packet.SYN == true {
			connection.ShakeHands()
			connection.State = SYN_RECV
		} else {
			connection.RST()
			connection.State = LISTEN
		}
	case SYN_RECV:
		connection.Comfirm(packet)
	case ESTABLISHED:
	case CLOSE_WAIT:
	case LAST_ACK:
	}
}



//
//	// 握手包TCP层
//	reqTCP := request.Layer(layers.LayerTypeTCP).(*layers.TCP)
//	// 握手包IP层
//	reqIP := request.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
//	// 握手包以太网层
//	reqETH := request.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)
//
//	// 最大TCP段长度
//	MSS := uint16(1469)
//	for _,v := range reqTCP.Options {
//		if v.OptionType == layers.TCPOptionKindMSS {
//			MSS = binary.BigEndian.Uint16(v.OptionData)
//		}
//	}
//
//	// 新建连接
//	conn := &Connection{
//		srcIP:		reqIP.DstIP,
//		srcPort:	reqTCP.DstPort,
//		srcMac:		reqETH.DstMAC,
//		srcSeq:		0,
//		dstIP:		reqIP.SrcIP,
//		dstPort:	reqTCP.SrcPort,
//		dstMac:		reqETH.SrcMAC,
//		dstSeq:		reqTCP.Seq,
//		dstMSS:		MSS,
//		dstWin:		reqTCP.Window,
//		Channel:	channel,
//		State: 		UNCONNECT,
//	}
//	// 进行TCP握手 发送SYN
//	conn.sendSYN()
//	conn.State = WAITSYNACK
//	return conn
//}

// 更新Connect数据
func (connection *Connection)Comfirm(rawPacket gopacket.Packet) {
	tcp := rawPacket.Layer(layers.LayerTypeTCP).(*layers.TCP)
	connection.srcSeq = tcp.Ack
	connection.dstSeq = tcp.Seq + uint32(len(tcp.Payload))
	connection.dstAck = tcp.Ack
	connection.dstWin = tcp.Window
}

func (connection *Connection)RST() {

}

func (connection *Connection)ShakeHands() {

}