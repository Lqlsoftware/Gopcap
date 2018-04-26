package http

import (
	"os"
	"regexp"
)

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func checkType(url string) bool {
	ct := getContentType(url)
	ok,_ := regexp.MatchString("text/|application/", ct)
	return ok
}