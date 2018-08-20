# Gopcap
A Multithreading HTTP web server based on pcap TCP layer.

## Quick start

```sh
# assume the following codes in main.go file
$ cat main.go
```

```go
package main

import (
    "github.com/Lqlsoftware/gopcap"
    "github.com/Lqlsoftware/gopcap/http"
)

func main() {
    gopcap.Bind("/", http.GET, handler)
    gopcap.Start(80) // serve on 80 port(http)
}

func handler(req *http.HttpRequest, rep *http.HttpResponse) {
    rep.Write("Hello World!\n")
}
```

```
# run main.go and server will start.
$ go run main.go
```
## Using
- Download and install it:
```sh
$ go get github.com/Lqlsoftware/Gopcap
```
- Import package in your code:
``` go
import "github.com/Lqlsoftware/gopcap"
```
- Write a handle function like:
``` go
func handler(req *http.HttpRequest, rep *http.HttpResponse) {
    rep.Write("Hello World!\n")
}
```
- Bind your handle function with an URL in your main function:
``` go
gopcap.Bind("/", http.GET, handler)
```
- Start server with port:
``` go
gopcap.Start(80)
```
- Run your project and enjoy!

- Put static html file in `root` folder (generate automaticly).

## Php7
- Download go-php library github.com/deuill/go-php
```sh
$ go get github.com/deuill/go-php
```
- Download Gopcap with tags php:
```sh
$ go get -tags php github.com/Lqlsoftware/Gopcap
```
or build Gopcap with tags php:
```sh
$ go build -tags php github.com/Lqlsoftware/Gopcap
```
- Enable Gopcap php in your main function:
``` go
gopcap.SetUsePhp()
```
- .php file will be excuted automaticaly.
