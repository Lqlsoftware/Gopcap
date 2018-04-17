package tcp

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TCP写入数据接口
func (conn *Connection)WriteData(data []byte) {
	conn.writeSlice(data)
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
	conn.srcSeq += uint32(len(data))
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
	fmt.Println("Send:\n",gopacket.NewPacket(buf.Bytes(),layers.LayerTypeEthernet,gopacket.DecodeOptions{}))
	check(err)
}