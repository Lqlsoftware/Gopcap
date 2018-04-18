package main

import (
	"gopcap"
	"gopcap/http"
)

func main() {
	gopcap.Bind("/", http.GET, handler)
	// 启动服务器
	gopcap.Start(8998)
}

func handler(req *http.HttpRequest, rep*http.HttpResponse) {
	rep.Write("Hello World!")
}