package gopcap

import (
	"log"
	"os"

	"github.com/Lqlsoftware/gopcap/http"
	"github.com/Lqlsoftware/gopcap/php"
	"github.com/google/gopacket/layers"
)

// 启动服务器
func Start(port layers.TCPPort)  {
	log.SetPrefix("[Gopcap] ")
	// server根目录root文件夹
	mkRoot()
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

func SetUsePhp() {
	php.UsePhp = true
}

// 创建根目录root文件夹
func mkRoot() {
	if !checkDirIsExist("root") {
		err := os.Mkdir("root", os.ModePerm)
		check(err)
	}
	if !checkDirIsExist("root/_temp") {
		err := os.Mkdir("root/_temp", os.ModePerm)
		check(err)
	}
}