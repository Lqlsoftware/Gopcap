package http

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"strings"
)

// Default GET method
// 		root/URL
func DefaultGETHandler(request *HttpRequest, response *HttpResponse) {
	dat, err := ioutil.ReadFile("root" + *request.url)
	if err != nil {
		response.stateCode = NotFound
		response.contents = []byte("<html>ERROR 404!</html>")
		return
	}

	// gzip 压缩
	if encoding,ok := (*request.header)["Accept-Encoding"];ok {
		encodes := strings.Split(encoding, ", ")
		// 检查浏览器是否支持gzip压缩
		for _,v := range encodes {
			// 支持gzip压缩
			if v == "gzip" {
				// 压缩数据
				var b bytes.Buffer
				w := gzip.NewWriter(&b)
				w.Write(dat)
				w.Flush()
				dat = b.Bytes()

				// 设置返回header 通知浏览器压缩格式
				(*response.header)["Content-Encoding"] = "gzip"
				break
			}
		}
	}

	// 设置Content-Type
	(*response.header)["Content-Type"] = getContentType(*request.url) + "; charset=utf-8"
	response.stateCode = OK
	response.contents = dat
}

// Default POST method
func DefaultPOSTHandler(request *HttpRequest, response *HttpResponse) {
	response.stateCode = OK
	response.contents = []byte("POST REQUEST: " + *request.url)
}

// Default HEAD method
func DefaultHEADHandler(request *HttpRequest, response *HttpResponse) {
	response.stateCode = OK
	response.contents = []byte{}
}