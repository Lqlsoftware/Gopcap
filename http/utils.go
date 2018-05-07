package http

import (
	"log"
	"os"
	"regexp"
	"time"
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

// 获取文件修改时间
func getFileModTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}