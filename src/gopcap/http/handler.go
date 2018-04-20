package http

import "fmt"

func Handler(rawPacket []byte) []byte {
	request,err := parserRequest(rawPacket)
	if err != nil {
		return []byte{}
	}
	fmt.Println(getmethodName(request.method),*request.url)
	response := request.generateResponse()
	// get router key:  url[0] ^= method -> url
	key := []byte(*request.url)
	key[0] ^= uint8(request.method)
	if handler,exist := routerMap[string(key)];exist {
		handler(request, response)
	} else {
		switch request.method {
		case GET:
			DefaultGETHandler(request, response)
		case POST:
			DefaultPOSTHandler(request, response)
		case HEAD:
			DefaultHEADHandler(request, response)
		default:
			DefaultGETHandler(request, response)
		}
	}
	return response.getBytes()
}