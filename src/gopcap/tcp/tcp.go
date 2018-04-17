package tcp

const tcpTimeout = 200

type State uint8
const (
	UNCONNECT	State = 0
	CONNECTED	State = 1
	SENDSYN		State = 2
	WAITSYNACK	State = 3
	SENDDATA	State = 4
	RECEIVEDATA	State = 5
	WAITACK		State = 6
	SENDFIN		State = 7
	RECEIVEFIN	State = 8
	WAITFINACK	State = 9
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}