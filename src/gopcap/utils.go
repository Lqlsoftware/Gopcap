package gopcap

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