package tcp

const tcpTimeout = 500000000

// TCP State
type State uint8
const (
	UNCONNECT	State = 0
	CONNECTED	State = 1
	WAITSYNACK	State = 3
	SENDDATA	State = 4
	RECEIVEDATA	State = 5
	WAITACK		State = 6
	SENDFIN		State = 7
	RECEIVEFIN	State = 8
	WAITFINACK	State = 9
)

// 错误处理
func check(err error) {
	if err != nil {
		panic(err)
	}
}