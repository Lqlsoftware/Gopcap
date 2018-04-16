package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const tcpTimeout = 200

type TcpState uint8
const (
	UNCONNECT	TcpState = 0
	CONNECTED	TcpState = 1
	WAITACK		TcpState = 2
	SENDDATA	TcpState = 3
	FIN			TcpState = 3
	WAITFINACK	TcpState = 4
)

type Connection struct {
	// 本机信息
	SrcIP		[]byte
	SrcPort		layers.TCPPort
	SrcMac		[]byte
	SrcSeq		uint32
	// 客户端信息
	DstIP		[]byte
	DstPort		layers.TCPPort
	DstMac		[]byte
	DstSeq		uint32
	Channel		chan gopacket.Packet
	State		TcpState
}

// TCP写入数据接口
func (conn *Connection)WriteData(data []byte) {
	if len(data) <= 1400 {
		conn.write(data)
	} else {
		conn.writeSlice(data)
	}
}

// 分段写入数据
func (conn *Connection)writeSlice(data []byte) {
	// 分片发送
	idx := 0
	for len(data) - idx > 1400 {
		buf := data[idx:idx + 1400]
		conn.write(buf)
		idx += 1400
	}
	buf := data[idx:]
	conn.write(buf)
}

// 写入小于1400字节的数据
func (conn *Connection)write(data []byte) {
	// ethernet
	ethLayer := layers.Ethernet{
		SrcMAC:		conn.SrcMac,
		DstMAC:		conn.DstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}
	// ip
	ipLayer := layers.IPv4{
		SrcIP:    	conn.SrcIP,
		DstIP:    	conn.DstIP,
		TTL:	  	64,
		Protocol: 	layers.IPProtocolTCP,
		Version:  	4,
		Flags:	  	2,
	}
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: 	conn.SrcPort,
		DstPort: 	conn.DstPort,
		ACK:	 	true,
		Ack:	 	conn.DstSeq,
		Seq:	 	conn.SrcSeq,
		Window:  	0xFFFF,
	}
	conn.SrcSeq += uint32(len(data))
	tcpLayer.Payload = data
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	buf := gopacket.NewSerializeBuffer()
	appLayer := gopacket.Payload(data)
	appLayer.SerializeTo(buf, gopacket.SerializeOptions{false,false})
	conn.writeRaw(buf, &ethLayer, &ipLayer, &tcpLayer)
}

// 封装TCP/IP/以太网包并发送
func (conn *Connection)writeRaw(buf gopacket.SerializeBuffer, ethLayer *layers.Ethernet, ipLayer *layers.IPv4, tcpLayer *layers.TCP) {
	tcpLayer.SerializeTo(buf, gopacket.SerializeOptions{true,true})
	ipLayer.SerializeTo(buf, gopacket.SerializeOptions{true,true})
	ethLayer.SerializeTo(buf, gopacket.SerializeOptions{false,true})
	err := handle.WritePacketData(buf.Bytes())
	check(err)
}

// 新建连接
func NewConnection(request gopacket.Packet) *Connection {
	reqTCP := request.Layer(layers.LayerTypeTCP).(*layers.TCP)
	reqIP := request.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	reqETH := request.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)
	// 建立端口channel
	channel := make(chan gopacket.Packet)
	// lock写入channel Map
	useMap<-true
	chMap[reqTCP.SrcPort] = channel
	<-useMap
	// 新建连接
	conn := &Connection{
		SrcIP:		reqIP.DstIP,
		SrcPort:	reqTCP.DstPort,
		SrcMac:		reqETH.DstMAC,
		SrcSeq:		0,
		DstIP:		reqIP.SrcIP,
		DstPort:	reqTCP.SrcPort,
		DstMac:		reqETH.SrcMAC,
		DstSeq:		reqTCP.Seq,
		Channel:	channel,
		State: 		UNCONNECT,
	}
	// 进行TCP握手
	conn.handShake()
	return conn
}

// 更新Connect数据
func (conn *Connection)Update(rawPacket gopacket.Packet) {
	tcp := rawPacket.Layer(layers.LayerTypeTCP).(*layers.TCP)
	conn.SrcSeq = tcp.Ack
	conn.DstSeq = tcp.Seq + uint32(len(tcp.Payload))
}

// 连接关闭
func (conn *Connection)Close() {
	// 从Map中删除
	useMap<-true
	delete(chMap, conn.DstPort)
	<-useMap
	// 关闭通道
	close(conn.Channel)
}

// TCP握手
func (conn *Connection)handShake() {
	// 发送SYN
	conn.sendSYN()
	// 等待client的ACK
	if conn.getACK() {
		conn.State = CONNECTED
	}
}

// 发送SYN
func (conn *Connection)sendSYN() {
	// ethernet
	ethLayer := layers.Ethernet{
		SrcMAC:		conn.SrcMac,
		DstMAC:		conn.DstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}
	// ip
	ipLayer := layers.IPv4{
		SrcIP:    	conn.SrcIP,
		DstIP:    	conn.DstIP,
		TTL:	  	64,
		Protocol: 	layers.IPProtocolTCP,
		Version:  	4,
		Flags:	  	2,
	}
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: 	conn.SrcPort,
		DstPort: 	conn.DstPort,
		SYN:     	true,
		ACK:	 	true,
		Ack:	 	conn.DstSeq + 1,
		Window:  	0xFFFF,
		Options:	[]layers.TCPOption{{layers.TCPOptionKindMSS,4,[]byte{5,189}}},
	}
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	buf := gopacket.NewSerializeBuffer()
	conn.writeRaw(buf, &ethLayer, &ipLayer, &tcpLayer)
}

// 发送ACK
func (conn *Connection)sendAck() {
	// ethernet
	ethLayer := layers.Ethernet{
		SrcMAC:		conn.SrcMac,
		DstMAC:		conn.DstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}
	// ip
	ipLayer := layers.IPv4{
		SrcIP:    	conn.SrcIP,
		DstIP:    	conn.DstIP,
		TTL:	  	64,
		Protocol: 	layers.IPProtocolTCP,
		Version:  	4,
		Flags:	  	2,
	}
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: 	conn.SrcPort,
		DstPort: 	conn.DstPort,
		ACK:	 	true,
		Ack:	 	conn.DstSeq,
		Seq:	 	conn.SrcSeq,
		Window:  	0xFFFF,
	}
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	buf := gopacket.NewSerializeBuffer()
	conn.writeRaw(buf, &ethLayer, &ipLayer, &tcpLayer)
}

// 发送FIN
func (conn *Connection)sendFin() {
	// ethernet
	ethLayer := layers.Ethernet{
		SrcMAC:		conn.SrcMac,
		DstMAC:		conn.DstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}
	// ip
	ipLayer := layers.IPv4{
		SrcIP:    	conn.SrcIP,
		DstIP:    	conn.DstIP,
		TTL:	  	64,
		Protocol: 	layers.IPProtocolTCP,
		Version:  	4,
		Flags:	  	2,
	}
	// tcp
	tcpLayer := layers.TCP{
		SrcPort: 	conn.SrcPort,
		DstPort: 	conn.DstPort,
		ACK:	 	true,
		FIN:		true,
		Ack:	 	conn.DstSeq,
		Seq:	 	conn.SrcSeq,
		Window:  	0xFFFF,
	}
	tcpLayer.SetNetworkLayerForChecksum(&ipLayer)
	buf := gopacket.NewSerializeBuffer()
	conn.writeRaw(buf, &ethLayer, &ipLayer, &tcpLayer)
}

// 等待ACK
func (conn *Connection)getACK() bool {
	time := 0
	for {
		select {
		case <-conn.Channel:
			return true
		}
		// 超时关闭连接
		time++
		if time > tcpTimeout {
			return false
		}
	}
}