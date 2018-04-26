package gopcap

import (
	"testing"

	"github.com/Lqlsoftware/gopcap/http"
)

func TestServer(t *testing.T) {
	// bind url router
	Bind("/helloWorld", http.GET, rootHandler)
	// start server
	Start(80)
}

// example handler
func rootHandler(req *http.HttpRequest, rep *http.HttpResponse) {
	rep.Write("Hello World!\n")
	for _,v := range req.GetAllParamKey() {
		rep.Write(v,":",req.GetParam(v),"\n")
	}
}