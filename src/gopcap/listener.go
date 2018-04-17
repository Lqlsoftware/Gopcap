package gopcap

import (
	"strconv"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"gopcap/handler"
	"fmt"
)

// 监听TCP端口
func listen(adapter *pcap.Interface, port layers.TCPPort) {
	in := openChannel(adapter, port)
	// 监听启动
	fmt.Print("Start listening Port:", port, "\n\n")
	for true {
		select {
		case packet := <-in:
			// tcp包
			handler.PacketHandler(packet)
		}
	}
}

func openChannel(adapter *pcap.Interface, port layers.TCPPort) chan gopacket.Packet {
	// 打开输入流
	handle,err := pcap.OpenLive(adapter.Name, 65535, true, pcap.BlockForever)
	defer handle.Close()
	check(err)
	// 设置过滤器
	err = handle.SetBPFFilter("tcp and dst port " + strconv.Itoa(int(port)))
	check(err)
	// 建立数据源
	src := gopacket.NewPacketSource(handle, handle.LinkType())
	checkNil(src)
	handler.SetConn(handle)
	return src.Packets()
}