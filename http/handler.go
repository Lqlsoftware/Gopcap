package http

import (
	"log"

	"github.com/Lqlsoftware/gopcap/php"
	"github.com/Lqlsoftware/gopcap/stream"
)

// HTTP包处理
func Handler(rawPacket []byte, phpPlugin *php.Plugin) (rep *stream.HttpStream, isKeepAlive bool) {
	// 转换TCP交付的包为HTTP-REQUEST
	request,err := parserRequest(rawPacket)
	if err != nil {
		return stream.NewFileStream(), false
	}

	// 判断是否keep-alive
	if (*request.header)["Connection"] == "keep-alive" {
		isKeepAlive = true
	}

	// 生成HTTP-RESPONSE
	response := request.generateResponse()

	// Get router key:  url[0] ^= method -> url
	key := []byte(*request.url)
	key[0] ^= uint8(request.method)

	// 查找URL路由表
	if handler,exist := routerMap[string(key)];exist {
		handler(request, response)
		response.stateCode = OK
	} else {
		switch request.method {
		case GET:
			DefaultGETHandler(request, response, phpPlugin)
		case POST:
			DefaultPOSTHandler(request, response)
		case HEAD:
			DefaultHEADHandler(request, response)
		default:
			DefaultGETHandler(request, response, phpPlugin)
		}
	}

	// 控制台log请求内容
	log.Println("|",response.stateCode,"|",getmethodName(request.method),*request.url)
	return response.getStream(), isKeepAlive
}