package http

import (
	"strings"
	"time"
)

// HTTP请求
type HttpRequest struct {
	// URL
	url			*string
	// HTTP-Version
	version		*string
	// HTTP-Header
	header 		*map[string]string
	// HTTP-Method
	method		HttpMethod
	// HTTP-Content
	contents	*string
	// Param
	param		*map[string]string
}

// 获取request的参数
func (req *HttpRequest)GetParam(key string) string {
	return (*req.param)[key]
}

// 获取request的所有参数名
func (req *HttpRequest)GetAllParamKey() []string {
	res := make([]string,0,len(*req.param))
	for v := range *req.param {
		res = append(res, v)
	}
	return res
}

// 根据REQUEST产生默认RESPONSE
func (req *HttpRequest)generateResponse() *HttpResponse {
	header := make(map[string]string)
	header["Server"] = "Gopcap"
	header["Date"] = time.Now().String()
	header["Accept-Ranges"] = "bytes"

	return &HttpResponse{
		header:		&header,
		version: 	req.version,
	}
}

// 处理REQUEST参数
func (req *HttpRequest)parseParameter() {
	parameter := make(map[string]string)
	switch req.method {
	case GET:
		s := strings.IndexByte(*req.url,'?')
		if s >= 0 {
			param := strings.Split((*req.url)[s + 1:], "&")
			*req.url = (*req.url)[:s]
			for _,v := range param {
				idx := strings.IndexByte(v,'=')
				key := v[:idx]
				value := v[idx + 1:]
				parameter[key] = value
			}
		}
	case POST:
		contentType := (*req.header)["Content-Type"]
		switch contentType {
		case "application/x-www-form-urlencoded":
			param := strings.Split(*req.contents, "&")
			for _,v := range param {
				idx := strings.IndexByte(v,'=')
				if idx < 0 {
					return
				}
				key := v[:idx]
				value := v[idx + 1:]
				parameter[key] = value
			}
		default:
			s := strings.Index(contentType,"boundary=")
			if s >= 0 {
				boundary := contentType[s + 9:] + "\r\n"
				if len(boundary) != 0 {
					param := strings.Split(*req.contents, boundary)
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

	}
	req.param = &parameter
}