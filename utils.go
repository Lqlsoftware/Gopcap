package gopcap

import "os"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkNil(v interface{}) {
	if v == nil {
		panic("nil interface")
	}
}

// 判断文件夹是否存在
func checkDirIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}