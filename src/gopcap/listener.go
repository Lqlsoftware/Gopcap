package gopcap

import (
	"log"
	"strconv"

	"gopcap/tcp"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// 监听TCP端口
func listen(adapter *pcap.Interface, port layers.TCPPort) {
	// 打开输入流
	handle,err := pcap.OpenLive(adapter.Name, 65535, true, pcap.BlockForever)
	tcp.SetSendChannel(handle)
	defer handle.Close()
	check(err)

	// 设置过滤器
	err = handle.SetBPFFilter("tcp and dst port " + strconv.Itoa(int(port)))
	check(err)

	// 建立数据源
	src := gopacket.NewPacketSource(handle, handle.LinkType())
	checkNil(src)
	in := src.Packets()

	// 监听启动
	log.Print("Port: ", port, "\n\n")
	for true {
		select {
		case packet := <-in:
			// tcp包
			tcp.PacketHandler(packet)
		}
	}
}