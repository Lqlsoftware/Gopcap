package http

import (
	"errors"
	"strings"
)

// 获取请求文件的对应ContentType
func getContentType(fileName string) string {
	idx := strings.IndexByte(fileName, '.')
	if idx < 0 {
		return "application/octet-stream"
	}
	ct := fileName[idx:]
	if res,ok := typeMap[ct];ok {
		return res
	} else {
		return "application/octet-stream"
	}
}

// 获取HTTP对应状态码的字符串
func getStateName(state HttpStateCode) string {
	switch state {
	case OK:
		return "OK"
	case PartialContent:
		return "Partial Content"
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

// 获取HTTP请求类型的字符串
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

// 将TCP交付的HTTP数据包进行解析 返回HTTPREQUEST
func parserRequest(raw []byte) (*HttpRequest,error) {
	// REQUEST
	// 		METHOD URL HTTP-VERSION
	idx,start := 0,0
	for idx < len(raw) && raw[idx] != 13 {
		idx++
	}
	first := strings.Split(string(raw[:idx])," ")
	if len(first) < 3 {
		return nil,errors.New("ERROR: UNKNOWN HTTP CONTENT")
	}
	// 中文URL UNESCAPE
	url, err := unescape(first[1])
	if err != nil {
		return nil,errors.New("ERROR: UNKNOWN HTTP URL ENCODE")
	}
	version := first[2]

	// REQUEST-HEADER
	// 		Connection: keep-alive
	idx += 2
	start = idx
	header := make(map[string]string)
	var key string
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

	// REQUEST-CONTENT
	var contents string
	if idx < len(raw) {
		contents = string(raw[idx:])
	}


	// Generate request
	request := &HttpRequest{
		method: 	HttpMethod(raw[0]),
		header: 	&header,
		contents: 	&contents,
		url:		&url,
		version:	&version,
	}

	// Parse request parameter
	request.parseParameter()

	return request, nil
}