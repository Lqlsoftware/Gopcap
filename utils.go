package gopcap

import (
	"net"
	"os"

	"github.com/google/gopacket/pcap"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkNil(v interface{}) {
	if v == nil {
		panic("nil interface")
	}
}

// 判断文件夹是否存在
func checkDirIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func getIPV4(adapter *pcap.Interface) net.IP {
	for _,v := range adapter.Addresses {
		if len(v.IP) == 4 {
			return v.IP
		}
	}
	return net.IPv4(0,0,0,0)
}