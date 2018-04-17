package handler

import (
	"gopcap/http"
	"gopcap/tcp"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var conn *pcap.Handle
func SetConn(handle *pcap.Handle) {
	conn = handle
}

func PacketHandler(packet gopacket.Packet) {
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
	tcpConn := tcp.NewConnection(conn, channel, synPacket)
	// 超时计时器
	timer := NewTimer(tcp.TcpTimeout)
	for {
		select {
		case request := <-*tcpConn.Channel:
			switch tcpConn.State {
			case tcp.UNCONNECT:
				tcpConn = tcp.NewConnection(conn, channel, request)
			case tcp.CONNECTED:
				tcpConn.Update(request)
				tcpConn.sendAck()
				http.HttpHandler(tcpConn,request)
				tcpConn.State = tcp.WAITACK
			case tcp.WAITACK:
				if request.TransportLayer().(*layers.TCP).Ack < tcpConn.SrcSeq {
					continue
				}
				tcpConn.Update(request)
				tcpConn.sendFin()
				tcpConn.State = tcp.FIN
			case tcp.FIN:
				tcpConn.Update(request)
				tcpConn.State = tcp.WAITFINACK
			case tcp.WAITFINACK:
				tcpConn.Update(request)
				tcpConn.DstSeq++
				tcpConn.sendAck()
				tcpConn.State = tcp.UNCONNECT
				timer.Reset()
			}
		}
		if tcpConn.State == tcp.UNCONNECT {
			if timer.Tick() {
				return
			}
		}
	}
}