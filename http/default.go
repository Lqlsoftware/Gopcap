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
	// 检查URL
	if *request.url == "/" {
		response.stateCode = OK
		response.contents = []byte(defaultIndex)
		return
	} else if strings.HasSuffix(*request.url,"/") {
		response.stateCode = Forbidden
		response.contents = []byte(page403)
		return
	}

	useGzip,useCache := false,false
	cachePath := "root/_temp" + *request.url
	filePath := "root" + *request.url

	// 检查支持 gzip 压缩
	if encoding,ok := (*request.header)["Accept-Encoding"];ok {
		encodes := strings.Split(encoding, ", ")
		// 检查浏览器是否支持gzip压缩
		for _,v := range encodes {
			// 支持gzip压缩
			if v == "gzip" {
				useGzip = true
				break
			}
		}
	}

	// 检查缓存是否有已压缩文件
	if useGzip && checkFileIsExist(cachePath) && checkFileIsExist(filePath) {
		// 检查缓存文件和新文件的修改时间
		if getFileModTime(cachePath) >= getFileModTime(filePath) {
			useCache = true
		}
	}

	var dat []byte
	var err error
	if useCache {
		// 设置返回header 通知浏览器压缩格式
		(*response.header)["Content-Encoding"] = "gzip"

		// 直接返回缓存数据
		dat,err = ioutil.ReadFile(cachePath)
	} else {
		// 读入文件
		dat, err = ioutil.ReadFile("root" + *request.url)
		if err != nil {
			response.stateCode = NotFound
			response.contents = []byte(page404)
			return
		}

		// gzip 压缩
		if useGzip {
			// 设置返回header 通知浏览器压缩格式
			(*response.headeyr)["Content-Encoding"] = "gzip"

			// 压缩数据
			var b bytes.Buffer
			w := gzip.NewWriter(&b)
			w.Write(dat)
			w.Flush()
			dat = b.Bytes()

			// 检查是否为静态text类文件
			if checkType(*request.url) {
				// 缓存
				ioutil.WriteFile(cachePath, dat, 0666)
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