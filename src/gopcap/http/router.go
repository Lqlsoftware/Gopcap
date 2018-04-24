package http

import (
	"errors"
)

// URL路由map
var routerMap = make(map[string]func(*HttpRequest,*HttpResponse))

func AddRouter(Url []byte, method HttpMethod, handler func(*HttpRequest,*HttpResponse)) error {
	if len(Url) == 0 {
		return errors.New("ERROR: URL is empty")
	}
	Url[0] ^= uint8(method)
	key := string(Url)
	if _,exist := routerMap[key];exist {
		return errors.New("ERROR: URL is already bind with handler")
	}
	routerMap[key] = handler
	return nil
}

func RemoveRouter(Url []byte, method HttpMethod) error {
	if len(Url) == 0 {
		return errors.New("ERROR: URL is empty")
	}
	Url[0] ^= uint8(method)
	key := string(Url)
	if _,exist := routerMap[key];!exist {
		return errors.New("ERROR: URL at " + getmethodName(method) + " is not bind")
	}
	delete(routerMap, key)
	return nil
}