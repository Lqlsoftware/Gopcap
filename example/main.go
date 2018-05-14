package main

import (
	"github.com/Lqlsoftware/gopcap"
	"github.com/Lqlsoftware/gopcap/http"
)

func main() {
	gopcap.SetUsePhp()
	// bind url router
	gopcap.Bind("/helloWorld", http.GET, rootHandler)
	// start server
	gopcap.Start(80)
}

// example handler
func rootHandler(req *http.HttpRequest, rep *http.HttpResponse) {
	rep.Write("Hello World!\n")
	for _,v := range req.GetAllParamKey() {
		rep.Write(v,":",req.GetParam(v),"\n")
	}
}