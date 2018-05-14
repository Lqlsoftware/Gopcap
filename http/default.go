package http

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/Lqlsoftware/gopcap/php"
)

// Default GET method
// 		root/URL
func DefaultGETHandler(request *HttpRequest, response *HttpResponse, phpPlugin *php.Plugin) {
	// 检查URL
	if *request.url == "/" {
		if !checkFileIsExist("root/index.html") {
			response.stateCode = OK
			response.contents = []byte(defaultIndex)
			return
		}
		*request.url = "/index.html"
	} else if strings.HasSuffix(*request.url,"/") {
		response.stateCode = Forbidden
		response.contents = []byte(page403)
		return
	} else if !checkFileIsExist("root" + *request.url) {
		response.stateCode = NotFound
		response.contents = []byte(page404)
		return
	} else if phpPlugin != nil && strings.HasSuffix(*request.url,".php") {
		var b bytes.Buffer
		w := io.Writer(&b)
		phpPlugin.SetPhpWriter(w)
		err := phpPlugin.Exec("root" + *request.url)
		if err != nil {
			response.stateCode = InternalServerError
			response.contents = []byte(defaultIndex)
			return
		}
		response.stateCode = OK
		response.contents = b.Bytes()
		(*response.header)["Content-Type"] = "text/html; charset=utf-8"
		return
	}

	cachePath := "root/_temp" + *request.url
	filePath := "root" + *request.url
	f,err := os.Open(filePath)
	check(err)
	defer f.Close()
	var dat []byte

	useGzip,useCache,useSlice := false	,false,false
	sliceSize,sliceStart,sliceEnd := getFileSize(f),int64(0),int64(0)



	// 检查断点续传和 gzip 压缩
	if slice,ok := (*request.header)["Range"];ok {
		useSlice = true
		start := strings.Index(slice, "bytes=") + 6
		if start < 0 || start > len(slice) {
			response.stateCode = Forbidden
			response.contents = []byte(page403)
			return
		}

		var err error
		for i := start;i < len(slice);i++ {
			if slice[i] == '-' {
				sliceStart,err = strconv.ParseInt(slice[start:i],10,64)
				if err != nil {
					response.stateCode = Forbidden
					response.contents = []byte(page403)
					return
				}

				// 默认全部
				if i + 1 == len(slice) {
					sliceEnd = sliceSize - 1
				} else {
					sliceEnd,err = strconv.ParseInt(slice[i + 1:],10,64)
					if err != nil || sliceSize <= sliceEnd {
						response.stateCode = Forbidden
						response.contents = []byte(page403)
						return
					}
				}
				break
			}
		}
		if etag,ok := (*request.header)["If-Range"];ok {
			if etag != strconv.FormatInt(getFileModTime(f),10) {
				useSlice = false
			}
		}
	} else if encoding,ok := (*request.header)["Accept-Encoding"];ok {
		encodes := strings.Split(encoding, ", ")
		// 检查浏览器是否支持gzip压缩
		for _,v := range encodes {
			// 支持gzip压缩
			if v == "gzip" {
				useGzip = true

				// 检查缓存是否有已压缩文件 缓存文件和新文件的修改时间
				if checkFileIsExist(cachePath) && getFileModTime(f) >= getFileModTime(f){
					useCache = true
				}
				break
			}
		}
	}

	if useCache {
		// 设置返回header 通知浏览器压缩格式
		(*response.header)["Content-Encoding"] = "gzip"

		// 直接返回缓存数据
		dat,err = ioutil.ReadFile(cachePath)
	} else if useSlice {
		response.SetHeader("Content-Range","bytes " + strconv.FormatInt(sliceStart,10) + "-" + strconv.FormatInt(sliceEnd,10) + "/" + strconv.FormatInt(sliceSize,10))

		buffer := make([]byte, sliceEnd - sliceStart + 1)
		_,err = f.ReadAt(buffer, sliceStart)
		check(err)

		if sliceStart == 0 {
			response.stateCode = OK
		} else {
			response.stateCode = PartialContent
		}
		response.contents = buffer
		return
	} else {
		// 读入文件
		dat, err = ioutil.ReadFile("root" + *request.url)
		check(err)

		// gzip 压缩
		if useGzip {
			// 设置返回header 通知浏览器压缩格式
			(*response.header)["Content-Encoding"] = "gzip"

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
		} else {
			// 文件最后修改时间 - 文件大小字节数转为16进制
			(*response.header)["ETag"] = strconv.FormatInt(getFileModTime(f),10)
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