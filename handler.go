package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func packetHandler(packet gopacket.Packet) {
	// 解析TCP报文
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	tcp := tcpLayer.(*layers.TCP)
	// 握手请求
	if _,ok := chMap[tcp.SrcPort];ok {
		// 发送至相应端口通道
		chMap[tcp.SrcPort]<- packet
	} else if tcp.SYN {
		// 建立新线程进行TCP连接
		go handleThread(packet)
	}
}

// 处理线程
func handleThread(connPacket gopacket.Packet) {
	conn := NewConnection(connPacket)
	defer conn.Close()
	if conn.State == UNCONNECT {
		return
	}
	time := 0
	for {
		select {
		case request := <-conn.Channel:
			switch conn.State {
			case UNCONNECT:
				conn = NewConnection(request)
			case CONNECTED:
				conn.Update(request)
				conn.sendAck()
				httpHandler(conn,request)
				conn.State = WAITACK
			case WAITACK:
				if request.TransportLayer().(*layers.TCP).Ack < conn.SrcSeq {
					continue
				}
				conn.Update(request)
				conn.sendFin()
				conn.State = FIN
			case FIN:
				conn.Update(request)
				conn.State = WAITFINACK
			case WAITFINACK:
				conn.Update(request)
				conn.DstSeq++
				conn.sendAck()
				conn.State = UNCONNECT
				time = 0
			}
		}
		if conn.State == UNCONNECT {
			time++
			if time > tcpTimeout {
				return
			}
		}
	}
}