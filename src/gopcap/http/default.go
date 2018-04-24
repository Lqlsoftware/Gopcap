package http

import "io/ioutil"

// Default GET method
// 		root/URL
func DefaultGETHandler(request *HttpRequest, response *HttpResponse) {
	dat, err := ioutil.ReadFile("root" + *request.url)
	if err != nil {
		response.stateCode = NotFound
		response.contents = []byte("<html>ERROR 404!</html>")
		return
	}
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