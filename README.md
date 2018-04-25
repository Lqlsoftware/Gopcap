# Gopcap
A Multithreading HTTP web server based on pcap TCP layer.

## Install

    go get github.com/Lqlsoftware/Gopcap

## Complie
Include Gopcap in your $GOPATH.
``` bash
export GOPATH=$GOPATH:(your go get dir)/github.com/Lqlsoftware/Gopcap
```
then

    go build gopcap
    go build main.go (Optional.) 

## Usage
Put static html to "./root/".

    go run main.go
    or
    ./main

---
## Develop
Import gopcap package:
``` go
import "github.com/Lqlsoftware/gopcap"
```
Write a handle function like:
``` go
func handler(req *http.HttpRequest, rep *http.HttpResponse) {
    rep.Write("Hello World!\n")
}
```
Bind your handle function with an URL:
``` go
gopcap.Bind("/", http.GET, handler)
```
Start server with port:
``` go
gopcap.Start(8998)
```
Enjoy!
