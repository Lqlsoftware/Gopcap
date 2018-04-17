package tcp

import (
	"gopcap/http"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var sendChannel *pcap.Handle
func SetSendChannel(channel *pcap.Handle) {
	sendChannel = channel
}

func PacketHandler(packet gopacket.Packet) {
	if packet == nil {
		return
	}
	// 解析TCP报文
	tcpLayer := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
	// 处理请求
	if _,ok := chMap[tcpLayer.SrcPort];ok {
		// 发送至相应端口通道
		chMap[tcpLayer.SrcPort]<- packet
	} else if tcpLayer.SYN {
		// 建立新线程进行TCP连接
		go handleThread(packet, tcpLayer.SrcPort)
	}
}

// 处理线程
func handleThread(synPacket gopacket.Packet, dstPort layers.TCPPort) {
	// 建立端口channel 写入Map
	channel := addChannel(dstPort)
	defer delChannel(dstPort)
	// 建立TCP连接
	tcpConn := NewConnection(channel, synPacket)
	// 超时计时器
	timer := NewTimer(tcpTimeout)
	for {
		select {
		case request := <-*tcpConn.Channel:
			switch tcpConn.State {
			case UNCONNECT:
				tcpConn = NewConnection(channel, request)
			case CONNECTED:
				tcpConn.Update(request)
				tcpConn.sendAck()
				response := http.HttpHandler(request)
				tcpConn.WriteData(response)
				tcpConn.State = WAITACK
			case WAITACK:
				if request.TransportLayer().(*layers.TCP).Ack < tcpConn.srcSeq {
					continue
				}
				tcpConn.Update(request)
				tcpConn.sendFin()
				tcpConn.State = SENDFIN
			case SENDFIN:
				tcpConn.Update(request)
				tcpConn.State = WAITFINACK
			case WAITFINACK:
				tcpConn.Update(request)
				tcpConn.dstSeq++
				tcpConn.sendAck()
				tcpConn.State = UNCONNECT
				timer.Reset()
			}
		}
		if tcpConn.State == UNCONNECT {
			if timer.Tick() {
				return
			}
		}
	}
}