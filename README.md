# Gopcap
A web server based on pcap TCP layer.

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
