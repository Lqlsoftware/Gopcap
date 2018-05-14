package http

import (
	"log"
	"os"
	"regexp"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

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
func getFileModTime(f *os.File) int64 {
	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

// 获取文件大小
func getFileSize(f *os.File) int64 {
	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return -1
	}

	return fi.Size()
}