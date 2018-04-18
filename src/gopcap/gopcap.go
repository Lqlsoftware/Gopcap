package gopcap

import (
	"github.com/google/gopacket/layers"
	"gopcap/http"
)

func Start(port layers.TCPPort)  {
	// 选择适配器
	adapter := getAdapter()
	// 开启TCP端口监听
	listen(adapter, port)
}

// dymatic bind
func Bind(Url string, method http.HttpMethod, handler func(*http.HttpRequest,*http.HttpResponse)) {
	err := http.AddRouter([]byte(Url), method, handler)
	check(err)
}

func DeBind(Url string, method http.HttpMethod) {
	err := http.RemoveRouter([]byte(Url), method)
	check(err)
}