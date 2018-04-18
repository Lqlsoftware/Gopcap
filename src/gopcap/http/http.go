package http

import (
	"github.com/google/gopacket"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
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

type HttpstateCode uint16
const (
	OK 						HttpstateCode = 200
	BadRequest				HttpstateCode = 400
	Unauthorized			HttpstateCode = 401
	Forbidden				HttpstateCode = 403
	NotFound				HttpstateCode = 404
	InternalServerError		HttpstateCode = 500
	ServerUnavailable		HttpstateCode = 503
)

var CRLF = []byte{13,10}
var SEP = []byte{58,32}

type HttpRequest struct {
	url			*string
	version		*string
	header 		*map[string]string
	method		HttpMethod
	contents	*string
	param		*map[string]string
}

type HttpResponse struct {
	header 		*map[string]string
	version		*string
	stateCode	HttpstateCode
	contents	[]byte
	ContentType	string
}

func (req *HttpRequest)GetParam(key string) string {
	return (*req.param)[key]
}

func (req *HttpRequest)GetAllParamKey() []string {
	res := make([]string,0,len(*req.param))
	for v := range *req.param {
		res = append(res, v)
	}
	return res
}


func (rep *HttpResponse)Write(Data ...interface{}) {
	for argNum, arg := range Data {
		if argNum > 0 {
			rep.contents = append(rep.contents,' ')
		}
		if arg != nil {
			rep.contents = append(rep.contents, []byte(reflect.ValueOf(arg).String())...)
		} else {
			rep.contents = append(rep.contents,[]byte("nil")...)
		}
	}
}

func (rep *HttpResponse)SetHeader(key string, value string) {
	(*rep.header)[key] = value
}

func generateResponse(req *HttpRequest) *HttpResponse {
	header := make(map[string]string)
	header["Server"] = "Gopcap"
	header["Date"] = time.Now().String()
	return &HttpResponse{
		header:		&header,
		version: 	req.version,
	}
}

func getStateName(state HttpstateCode) string {
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

func getmethodName(method HttpMethod) string {
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

func (rep *HttpResponse)getBytes() []byte {
	length := 38 + len(rep.contents)
	for key,value := range *rep.header {
		length += len(key) + len(value) + 4
	}
	buf := make([]byte, 0, length)
	buf = append(buf, []byte(*rep.version)...)
	buf = append(buf, 32)
	buf = append(buf, []byte(strconv.Itoa(int(rep.stateCode)))...)
	buf = append(buf, 32)
	buf = append(buf, []byte(getStateName(rep.stateCode))...)
	buf = append(buf, CRLF...)
	// header
	for key,value := range *rep.header {
		buf = append(buf, []byte(key)...)
		buf = append(buf, SEP...)
		buf = append(buf, []byte(value)...)
		buf = append(buf, CRLF...)
	}
	buf = append(buf, CRLF...)
	// content
	buf = append(buf, rep.contents...)
	return buf
}

func HttpHandler(rawPacket gopacket.Packet) []byte {
	if rawPacket.ApplicationLayer() == nil {
		return []byte{}
	}
	request,err := parser(rawPacket.ApplicationLayer().Payload())
	if err != nil {
		return []byte{}
	}
	fmt.Println(getmethodName(request.method),*request.url)
	response := generateResponse(request)
	key := []byte(*request.url)
	key[0] ^= uint8(request.method)
	if handler,exist := routerMap[string(key)];exist {
		handler(request, response)
	} else {
		switch request.method {
		case GET:
			GETHandler(request, response)
		case POST:
			POSTHandler(request, response)
		case HEAD:
			HEADHandler(request, response)
		default:
			GETHandler(request, response)
		}
	}
	return response.getBytes()
}

func parser(raw []byte) (*HttpRequest,error) {
	header := make(map[string]string)
	idx,start := 0,0
	key := ""
	for idx < len(raw) && raw[idx] != 13 {
		idx++
	}
	first := strings.Split(string(raw[:idx])," ")
	if len(first) < 3 {
		return nil,errors.New("ERROR: UNKNOWN HTTP CONTENT")
	}
	url := first[1]
	version := first[2]
	idx += 2
	start = idx
	for idx < len(raw) {
		if v := raw[idx];v == 13 {
			if idx <= start {
				break
			}
			header[key] = string(raw[start:idx])
			idx += 2
			start = idx
		} else if v == 58 && raw[idx + 1] == 32 {
			key = string(raw[start:idx])
			idx += 2
			start = idx
		} else {
			idx++
		}
	}
	contents := string(raw[idx:])
	parameter := make(map[string]string)
	if HttpMethod(raw[0]) == GET {
		s := strings.IndexByte(url,'?')
		if s >= 0 {
			param := strings.Split(url[s + 1:], "&")
			url = url[:s]
			for _,v := range param {
				idx := strings.IndexByte(v,'=')
				key := v[:idx]
				value := v[idx + 1:]
				parameter[key] = value
			}
		}
	} else if HttpMethod(raw[0]) == POST {
		contentType := header["content-type"]
		s := strings.Index(contentType,"boundary=")
		if s >= 0 {
			boundary := contentType[s + 9:] + "\r\n"
			if len(boundary) != 0 {
				param := strings.Split(contents, boundary)
				for _,v := range param {
					start := strings.Index(v,"name=\"")
					end := strings.Index(v,"\"\r\n\r\n")
					if start >= 0 && end >= 0 {
						key := v[start + 6:end]
						value := v[end + 5:]
						parameter[key] = value
					}
				}
			}
		}
	}
	request := &HttpRequest{
		method: 	HttpMethod(raw[0]),
		header: 	&header,
		contents: 	&contents,
		url:		&url,
		version:	&version,
		param:		&parameter,
	}
	return request, nil
}

func GETHandler(request *HttpRequest, response *HttpResponse) {
	dat, err := ioutil.ReadFile("root" + *request.url)
	(*response.header)["Content-Type"] = "text/html; charset=utf-8"
	if err != nil {
		response.stateCode = NotFound
		response.contents = []byte("<html>ERROR 404!</html>")
		return
	}
	response.stateCode = OK
	response.contents = dat
}

func POSTHandler(request *HttpRequest, response *HttpResponse) {
	response.stateCode = OK
	response.contents = []byte("POST REQUEST: " + *request.url)
}

func HEADHandler(request *HttpRequest, response *HttpResponse) {
	response.stateCode = OK
	response.contents = []byte{}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}