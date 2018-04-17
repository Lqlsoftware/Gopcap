package tcp

type Timer struct {
	current uint32
	max 	uint32
}

func NewTimer(max uint32) Timer {
	return Timer{
		current:	0,
		max:		max,
	}
}

func (timer *Timer)Reset() {
	timer.current = 0
}

func (timer *Timer)Tick() bool {
	timer.current++
	if timer.current >= timer.max {
		return true
	}
	return false
}