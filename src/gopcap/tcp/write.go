package tcp

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TCP写入数据接口
func (conn *Connection)WriteData(data []byte, startSeq uint32) {
	window := conn.dstWin
	start := conn.dstAck - startSeq
	conn.writeSlice(data, start, window)
}

// 分段写入数据
func (conn *Connection)writeSlice(data []byte, start uint32, window uint16) {
	// 分片发送
	end := start + uint32(window)
	length := uint32(len(data))
	if end > length {
		end = length
	}
	curr := start + 1400
	for end - start > 0 {
		if curr > end {
			curr = end
		}
		conn.write(data[start:curr])
		conn.srcSeq += uint32(curr - start)
		start = curr
		curr += 1400
	}
}

func (conn *Connection)Rewrite(data []byte, startSeq uint32) {
	idx := len(data) - int(conn.srcSeq - conn.dstAck)
	if idx < 0 || idx >= len(data) {
		return
	}
	conn.WriteData(data, startSeq)
}
// 写入小于1400字节的数据
func (conn *Connection)write(data []byte) {
	// upper layer
	ethLayer, ipLayer := conn.getUpperLayers()
	// tcp
	tcpLayer := conn.getLayers()
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	buf := gopacket.NewSerializeBuffer()
	appLayer := gopacket.Payload(data)
	err := appLayer.SerializeTo(buf, gopacket.SerializeOptions{false,false})
	check(err)
	conn.writeRaw(buf, ethLayer, ipLayer, tcpLayer)
}

// 封装TCP/IP/以太网包并发送
func (conn *Connection)writeRaw(buf gopacket.SerializeBuffer, ethLayer *layers.Ethernet, ipLayer *layers.IPv4, tcpLayer *layers.TCP) {
	err := tcpLayer.SerializeTo(buf, gopacket.SerializeOptions{true,true})
	check(err)
	err = ipLayer.SerializeTo(buf, gopacket.SerializeOptions{true,true})
	check(err)
	err = ethLayer.SerializeTo(buf, gopacket.SerializeOptions{false,true})
	check(err)
	err = sendChannel.WritePacketData(buf.Bytes())
	check(err)
}