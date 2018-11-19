package tcp

import (
	"github.com/google/gopacket/pcap"
)

type Sender struct {
	Pipe *pcap.Handle
}

func (sender Sender)Init(handle *pcap.Handle) {
	sender.Pipe = handle
}
