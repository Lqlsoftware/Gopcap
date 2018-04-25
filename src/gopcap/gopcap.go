package gopcap

import (
	"log"

	"gopcap/http"

	"github.com/google/gopacket/layers"
)

// 启动服务器
func Start(port layers.TCPPort)  {
	log.SetPrefix("[Gopcap] ")
	// 选择适配器
	adapter := getAdapter()
	// 开启TCP端口监听
	listen(adapter, port)
}

// 绑定URL
func Bind(Url string, method http.HttpMethod, handler func(*http.HttpRequest,*http.HttpResponse)) {
	err := http.AddRouter([]byte(Url), method, handler)
	check(err)
}

// 解绑URL
func DeBind(Url string, method http.HttpMethod) {
	err := http.RemoveRouter([]byte(Url), method)
	check(err)
}