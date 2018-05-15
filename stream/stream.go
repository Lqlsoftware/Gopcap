package stream

import (
	"bytes"
	"fmt"
	"os"
)

const bufferSize  = 1 << 16

type HttpStream struct {
	// file descriptor
	fd 		*os.File
	// file read offset
	start	uint32
	// file len
	flen	uint32
	// header content
	hc 		*[]byte
	// header length
	// length
	length	uint32
	// buffer
	content	*bytes.Buffer
	buffer  *bytes.Buffer
	// 0 - not init    1 - file    2 - []byte
	flag	uint8
}

func NewFileStream() *HttpStream {
	return &HttpStream{
		fd:		nil,
		hc:		nil,
		start:	0,
		length: 0,
		content:bytes.NewBuffer([]byte{}),
		buffer: nil,
		flag:	2,
	}
}

func (hs *HttpStream)Write(data []byte) {
	hs.content.Write(data)
}

func (hs *HttpStream)WriteString(data string) {
	hs.content.WriteString(data)
}

func (hs *HttpStream)Len() uint32 {
	return hs.length
}

func (hs *HttpStream)UpdateLen() {
	if hs.flag == 2 {
		hs.flen = uint32(hs.content.Len())
	}
	hs.length = uint32(len(*hs.hc)) + hs.flen
	hs.buffer = bytes.NewBuffer(make([]byte, 0,bufferSize))
}

func (hs *HttpStream)SetRawHeader(h []byte) {
	hs.hc = &h
}

func (hs *HttpStream)SetFileDescriptor(f *os.File, start, length uint32) {
	hs.fd = f
	hs.flag = 1
	hs.start = start
	hs.flen = length
}

func (hs *HttpStream)Output(start, end uint32) []byte {
	// 重置Buffer
	hs.buffer.Reset()

	// 比头部长度小 输出头部
	hcLen := uint32(len(*hs.hc))
	if start < hcLen {
		if end > hcLen {
			hs.buffer.Write((*hs.hc)[start:hcLen])
			start = hcLen
		} else {
			hs.buffer.Write((*hs.hc)[start:end])
			return hs.buffer.Bytes()
		}
	}
	end = end - hcLen + hs.start
	start = min(start - hcLen + hs.start, end)

	// 根据流类型输出
	if hs.flag == 1 {
		// file
		data := make([]byte, end - start)
		hs.fd.ReadAt(data, int64(start))
		hs.buffer.Write(data)
	} else if hs.flag == 2 {
		// []byte
		fmt.Println(start,end,hs.content.Len())
		hs.buffer.Write(hs.content.Bytes()[start:end])
	}
	return hs.buffer.Bytes()
}

func (hs *HttpStream)GetContentLen() int64 {
	if hs.flag == 2 {
		return int64(hs.content.Len())
	} else {
		return int64(hs.flen)
	}
}

func (hs *HttpStream)Close() {
	if hs.fd != nil {
		hs.fd.Close()
	}
}

func min(a, b uint32) uint32 {
	if a <= b {
		return a
	}
	return b
}