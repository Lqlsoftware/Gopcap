package handler

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var chMap = make(map[layers.TCPPort]chan gopacket.Packet)
var useMap = make(chan bool,1)

// 创建并加入Map
func addChannel(key layers.TCPPort) *chan gopacket.Packet {
	channel := make(chan gopacket.Packet)
	useMap<-true
	chMap[key] = channel
	<-useMap
	return &channel
}

// 从Map中删除
func delChannel(key layers.TCPPort) {
	channel := chMap[key]
	useMap<-true
	delete(chMap, key)
	<-useMap
	close(channel)
}