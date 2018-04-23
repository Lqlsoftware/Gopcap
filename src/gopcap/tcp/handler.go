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
	var response []byte
	var input []byte
	var startSeq uint32
	for {
		select {
		case request := <-*tcpConn.Channel:
			tcp := request.TransportLayer().(*layers.TCP)
			if tcp.Ack < tcpConn.srcSeq {
				tcpConn.dstAck = tcp.Ack
				continue
			}
			tcpConn.Update(request)
			switch tcpConn.State {
			case UNCONNECT:
				if tcp.SYN {
					tcpConn = NewConnection(channel, request)
				} else if tcp.FIN {
					tcpConn.sendFin()
					tcpConn.State = SENDFIN
				}
			case WAITSYNACK:
				tcpConn.State = CONNECTED
			case CONNECTED:
				tcpConn.sendAck()
				if request.ApplicationLayer() == nil {
					continue
				} else if len(tcp.Payload) == int(tcpConn.dstMSS) {
					// 请求长度超过MSS
					input = append(input, tcp.Payload...)
					continue
				}
				input = append(input, tcp.Payload...)
				response = http.Handler(input)
				startSeq = tcpConn.srcSeq
				tcpConn.WriteData(response, startSeq)
				tcpConn.State = SENDDATA
				timer.Reset()
			case SENDDATA:
				if tcpConn.dstAck >= startSeq + uint32(len(response)) {
					tcpConn.sendFin()
					tcpConn.State = SENDFIN
				} else {
					tcpConn.WriteData(response, startSeq)
				}
			case SENDFIN:
				tcpConn.State = WAITFINACK
			case WAITFINACK:
				if tcp.FIN {
					tcpConn.dstSeq++
					tcpConn.sendAck()
					tcpConn.State = UNCONNECT
					timer.Reset()
				}
			}
		}
		if tcpConn.State == UNCONNECT {
			// 超时关闭连接
			if timer.Tick() {
				return
			}
		} else if tcpConn.State == SENDDATA {
			// 超时重传
			if timer.Tick() {
				tcpConn.Rewrite(response, startSeq)
				timer.Reset()
			}
		}
	}
}