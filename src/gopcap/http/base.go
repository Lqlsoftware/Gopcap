package http

import (
	"errors"
	"strings"
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

func parserRequest(raw []byte) (*HttpRequest,error) {
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
	request := &HttpRequest{
		method: 	HttpMethod(raw[0]),
		header: 	&header,
		contents: 	&contents,
		url:		&url,
		version:	&version,
	}
	request.parseParameter()
	return request, nil
}