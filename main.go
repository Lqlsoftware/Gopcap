package main

import (
	"gopcap"
	"gopcap/http"
)

func main() {
	// bind url router
	gopcap.Bind("/", http.POST, rootHandler)
	// start server
	gopcap.Start(8998)
}

func rootHandler(req *http.HttpRequest, rep *http.HttpResponse) {
	rep.Write("Hello World!\n")
	for _,v := range req.GetAllParamKey() {
		rep.Write(v,":",req.GetParam(v),"\n")
	}
}