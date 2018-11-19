package tcp

import (
	"sync"

	"github.com/google/gopacket/layers"
)


const PortMapSize = 1 << 16

type Receiver struct {
	// TODO replace with concurrentMap
	PortMap 		map[layers.TCPPort]chan *layers.TCP
	PortMapLocker	*sync.RWMutex
}

func (receiver Receiver)Init() {
	receiver.PortMap = make(map[layers.TCPPort]chan *layers.TCP, PortMapSize)
	receiver.PortMapLocker = new(sync.RWMutex)
}

func (receiver Receiver)AddPort(port layers.TCPPort) *chan *layers.TCP {
	channel := make(chan *layers.TCP)
	receiver.PortMapLocker.Lock()
	receiver.PortMap[port] = channel
	receiver.PortMapLocker.Unlock()
	return &channel
}

func (receiver Receiver)RemovePort(port layers.TCPPort) {
	receiver.PortMapLocker.Lock()
	channel := receiver.PortMap[port]
	delete(receiver.PortMap, port)
	receiver.PortMapLocker.Unlock()
	close(channel)
}