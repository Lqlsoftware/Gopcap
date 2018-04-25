package tcp

import (
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var chMap = make(map[layers.TCPPort]chan gopacket.Packet)
var useMap *sync.RWMutex

// 创建并加入Map
func addChannel(key layers.TCPPort) *chan gopacket.Packet {
	channel := make(chan gopacket.Packet)
	useMap.Lock()
	chMap[key] = channel
	useMap.Unlock()
	return &channel
}

// 从Map中删除
func delChannel(key layers.TCPPort) {
	channel := chMap[key]
	useMap.Lock()
	delete(chMap, key)
	useMap.Unlock()
	close(channel)
}