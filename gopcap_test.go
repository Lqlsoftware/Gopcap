package gopcap

import (
	"github.com/Lqlsoftware/Gopcap/http"
)

func TestStart() {
	// bind url router
	Bind("/helloWorld", http.GET, rootHandler)
	// start server
	Start(80)
}

func rootHandler(req *http.HttpRequest, rep *http.HttpResponse) {
	rep.Write("Hello World!\n")
	for _,v := range req.GetAllParamKey() {
		rep.Write(v,":",req.GetParam(v),"\n")
	}
}