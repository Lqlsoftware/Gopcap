package gopcap

import (
	"github.com/google/gopacket/layers"
)

func Start(port layers.TCPPort)  {
	// 选择适配器
	adapter := getAdapter()
	// 开启TCP端口监听
	listen(adapter, port)
}