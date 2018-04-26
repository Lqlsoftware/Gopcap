package tcp

import (
	"github.com/Lqlsoftware/gopcap/http"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// 发送通道
var sendChannel *pcap.Handle
func SetSendChannel(channel *pcap.Handle) {
	sendChannel = channel
}

// TCP包处理
func PacketHandler(packet gopacket.Packet) {
	if packet == nil {
		return
	}

	// 解析TCP报文
	tcpLayer := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)

	// 处理请求 端口通道存在
	useMap.RLock()
	ch,ok := chMap[tcpLayer.SrcPort]
	useMap.RUnlock()
	if ok {
		// 发送至相应端口通道
		ch<- packet
	} else if tcpLayer.SYN {
		// 建立新线程进行TCP连接
		go handleThread(packet, tcpLayer.SrcPort)
	}
}

// 包处理线程
func handleThread(synPacket gopacket.Packet, dstPort layers.TCPPort) {
	// 建立端口channel 写入Map
	channel := addChannel(dstPort)
	defer delChannel(dstPort)

	// 建立TCP连接 开始握手
	tcpConn := NewConnection(channel, synPacket)

	// 超时计时器
	timer := NewTimer(tcpTimeout)

	// 状态变量设置
	var response, input []byte
	var startSeq, last, previous = uint32(0), uint32(0), uint32(0)
	var isKeepAlive bool

	// 处理后续TCP包
	for {
		select {
		case request := <-*tcpConn.Channel:
			// 解析TCP层
			tcp := request.TransportLayer().(*layers.TCP)

			if tcp.RST {
				// RST 重置连接
				tcpConn.Update(request)
				tcpConn.State = UNCONNECT
				tcpConn.sendAck()
				continue
			} else if tcp.FIN {
				// FIN 结束连接
				tcpConn.Update(request)
				tcpConn.dstSeq++
				tcpConn.sendAck()
				tcpConn.sendFin()
				tcpConn.State = SENDFIN
				continue
			} else if tcp.Ack < tcpConn.srcSeq {
				// 忽略连续重传请求
				if tcp.Ack == last && last == previous {
					continue
				}
				// 小于当前发送序号 且是第二次ACK
				if tcp.Ack == last {
					tcpConn.srcSeq = tcp.Ack
				}
				previous = last
				last = tcp.Ack
				tcpConn.dstAck = tcp.Ack
				tcpConn.dstWin = tcp.Window
				continue
			}

			// 更新连接
			tcpConn.Update(request)

			// 根据连接状态处理包
			switch tcpConn.State {
			// 未连接
			case UNCONNECT:
				if tcp.SYN {
					tcpConn = NewConnection(channel, request)
				} else if tcp.FIN {
					tcpConn.sendFin()
					tcpConn.State = SENDFIN
				}

			// 等待握手ACK
			case WAITSYNACK:
				tcpConn.State = CONNECTED

			// 已连接 / keep-alive
			case CONNECTED:
				// 返回ACK
				tcpConn.sendAck()

				if request.ApplicationLayer() == nil {
					// keep-alive 心跳包
					tcpConn.sendAck()
					continue
				} else if len(tcp.Payload) == int(tcpConn.dstMSS) {
					// 请求长度超过MSS 等待下个包
					input = append(input, tcp.Payload...)
					continue
				}
				input = append(input, tcp.Payload...)

				// 交由HTTP处理
				response,isKeepAlive = http.Handler(input)

				// 发送response
				input = nil
				startSeq = tcpConn.srcSeq
				tcpConn.WriteWindow(response, startSeq)
				tcpConn.State = SENDDATA
				timer.Reset()

			//	发送数据
			case SENDDATA:
				// 接收到最后序列的ACK 数据传输完成
				if tcpConn.dstAck >= startSeq + uint32(len(response)) {
					response = nil
					if isKeepAlive {
						// 保持连接
						tcpConn.State = CONNECTED
					} else {
						// 结束连接
						tcpConn.sendFin()
						tcpConn.State = SENDFIN
					}
				} else {
					// 继续发送数据
					tcpConn.WriteWindow(response, startSeq)
					timer.Reset()
				}

			// 发送FIN
			case SENDFIN:
				tcpConn.State = WAITFINACK

			// 等待FIN的ACK
			case WAITFINACK:
				if tcp.FIN {
					tcpConn.dstSeq++
					tcpConn.sendAck()
					// 重置连接 等待端口建立下个连接
					tcpConn.State = UNCONNECT
					timer.Reset()
				}
			}
		}

		// 计时器
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
		} else if tcpConn.State == CONNECTED {
			// keep-alive 超时关闭连接
			if timer.Tick() {
				tcpConn.sendFin()
				tcpConn.State = SENDFIN
				timer.Reset()
			}
		}
	}
}