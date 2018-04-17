package http

import (
	"github.com/google/gopacket"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type HttpMethod uint8
const (
	GET 	HttpMethod = 71
	POST	HttpMethod = 80
	HEAD	HttpMethod = 72
)

type HttpStateCode uint16
const (
	OK 						HttpStateCode = 200
	BadRequest				HttpStateCode = 400
	Unauthorized			HttpStateCode = 401
	Forbidden				HttpStateCode = 403
	NotFound				HttpStateCode = 404
	InternalServerError		HttpStateCode = 500
	ServerUnavailable		HttpStateCode = 503
)

var CRLF = []byte{13,10}
var SEP = []byte{58,32}

type httpRequest struct {
	Url			*string
	Version		*string
	Header 		*map[string]string
	Method		HttpMethod
	Contents	*string
}

type httpResponse struct {
	Header 		*map[string]string
	Version		*string
	StateCode	HttpStateCode
	Contents	*[]byte
	ContentType	string
}

func generateResponse(req *httpRequest) *httpResponse {
	header := make(map[string]string)
	header["Server"] = "Gopcap"
	header["Date"] = time.Now().String()
	return &httpResponse{
		Header:		&header,
		Version: 	req.Version,
	}
}

func getStateName(state HttpStateCode) string {
	switch state {
	case OK:
		return "OK"
	case BadRequest:
		return "Bad Request"
	case Unauthorized:
		return "Unauthorized"
	case InternalServerError:
		return "Internal Server Error"
	case ServerUnavailable:
		return "Server Unavailable"
	default:
		return ""
	}
}

func getMethodName(method HttpMethod) string {
	switch method {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case HEAD:
		return "HEAD"
	default:
		return ""
	}
}

func (rep *httpResponse)getBytes() []byte {
	(*rep.Header)["Content-Type"] = "text/html; charset=utf-8"
	length := 38// + len(*rep.Contents)
	for key,value := range *rep.Header {
		length += len(key) + len(value) + 4
	}
	buf := make([]byte, 0, length)
	buf = append(buf, []byte(*rep.Version)...)
	buf = append(buf, 32)
	buf = append(buf, []byte(strconv.Itoa(int(rep.StateCode)))...)
	buf = append(buf, 32)
	buf = append(buf, []byte(getStateName(rep.StateCode))...)
	buf = append(buf, CRLF...)
	// header
	for key,value := range *rep.Header {
		buf = append(buf, []byte(key)...)
		buf = append(buf, SEP...)
		buf = append(buf, []byte(value)...)
		buf = append(buf, CRLF...)
	}
	buf = append(buf, CRLF...)
	// content
	buf = append(buf, *rep.Contents...)
	return buf
}

func HttpHandler(rawPacket gopacket.Packet) []byte {
	if rawPacket.ApplicationLayer() == nil {
		return []byte{}
	}
	request := parser(rawPacket.ApplicationLayer().Payload())
	fmt.Println(getMethodName(request.Method),*request.Url)
	response := generateResponse(request)
	switch request.Method {
	case GET:
		GETHandler(request, response)
	case POST:
		POSTHandler(request, response)
	case HEAD:
		HEADHandler(request, response)
	default:
		GETHandler(request, response)
	}
	return response.getBytes()
	//conn.WriteData(response.getBytes())
}

func parser(raw []byte) *httpRequest {
	header := make(map[string]string)
	idx,start := 0,0
	key := ""
	for idx < len(raw) && raw[idx] != 13 {
		idx++
	}
	first := strings.Split(string(raw[:idx])," ")
	idx += 2
	start = idx
	for idx < len(raw) {
		if v := raw[idx];v == 13 {
			if idx <= start {
				break
			}
			idx += 2
			header[key] = string(raw[start:idx])
			start = idx
		} else if v == 58 && raw[idx + 1] == 32 {
			idx += 2
			key = string(raw[start:idx])
			start = idx
		} else {
			idx++
		}
	}
	contents := string(raw[idx:])
	request := &httpRequest{
		Method: 	HttpMethod(raw[0]),
		Header: 	&header,
		Contents: 	&contents,
		Url:		&first[1],
		Version:	&first[2],
	}
	return request
}

func GETHandler(request *httpRequest, response *httpResponse) {
	dat, err := ioutil.ReadFile("root" + *request.Url)
	if err != nil {
		response.StateCode = NotFound
		msg := []byte("<html>ERROR 404!</html>")
		response.Contents = &msg
		return
	}
	response.StateCode = OK
	response.Contents = &dat
}

func POSTHandler(request *httpRequest, response *httpResponse) {
	response.StateCode = OK
	msg := []byte("POST REQUEST: " + *request.Url)
	response.Contents = &msg
}

func HEADHandler(request *httpRequest, response *httpResponse) {
	var msg []byte
	response.StateCode = OK
	response.Contents = &msg
}