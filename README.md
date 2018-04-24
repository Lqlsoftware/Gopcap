# Gopcap
A HTTP web server based on pcap TCP layer.

## Install:

    go get github.com/Lqlsoftware/Gopcap

## Complie:
Include Gopcap in your $GOPATH.
```
export GOPATH=$GOPATH:(your go get dir)/github.com/Lqlsoftware/Gopcap
```
then

    go build gopcap
    go build main.go (Optional.) 

## Usage:
Put static html to "./root/".

    go run main.go
    or
    ./main

---
## Develop
Import gopcap package:
```
import "github.com/Lqlsoftware/gopcap"
```
Write a handle function like:
```
func handler(req *http.HttpRequest, rep *http.HttpResponse) {
    rep.Write("Hello World!\n")
}
```
Bind your handle function with an URL:
```
gopcap.Bind("/", http.GET, handler)
```
Start server with port:
```
gopcap.Start(8998)
```
Enjoy!
