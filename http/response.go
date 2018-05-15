package http

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/Lqlsoftware/gopcap/stream"
)

// "/r/n"
var CRLF = []byte{13,10}

// ": "
var SEP = []byte{58,32}

// HTTP返回
type HttpResponse struct {
	// HTTP-Header
	header 		*map[string]string
	// HTTP-Version
	version		*string
	// HTTP-State
	stateCode	HttpStateCode
	// HTTP-Content
	contents	*stream.HttpStream
}

// Response Content 写入接口
func (rep *HttpResponse)Write(Data ...interface{}) {
	for argNum, arg := range Data {
		if argNum > 0 {
			rep.contents.WriteString(" ")
		}
		if arg != nil {
			rep.contents.WriteString(reflect.ValueOf(arg).String())
		} else {
			rep.contents.WriteString("nil")
		}
	}
}

// 设置response的首部
func (rep *HttpResponse)SetHeader(key string, value string) {
	(*rep.header)[key] = value
}

// response变成HttpStream
func (rep *HttpResponse)getStream() *stream.HttpStream {
	// 设置默认头部
	(*rep.header)["Content-Length"] = strconv.FormatInt(rep.contents.GetContentLen(),10)

	// 计算byte总共长度 防止append申请内存拷贝
	length := 38
	for key,value := range *rep.header {
		length += len(key) + len(value)
	}

	// 申请固定capacity的内存
	buffer := bytes.NewBuffer(make([]byte, 0, length))
	buffer.WriteString(*rep.version)
	buffer.WriteByte(32)
	buffer.WriteString(strconv.Itoa(int(rep.stateCode)))
	buffer.WriteByte(32)
	buffer.WriteString(getStateName(rep.stateCode))
	buffer.Write(CRLF)

	// header
	for key,value := range *rep.header {
		buffer.WriteString(key)
		buffer.Write(SEP)
		buffer.WriteString(value)
		buffer.Write(CRLF)
	}
	buffer.Write(CRLF)

	rep.contents.SetRawHeader(buffer.Bytes())
	rep.contents.UpdateLen()
	return rep.contents
}