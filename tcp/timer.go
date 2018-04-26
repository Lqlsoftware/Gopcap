package tcp

// 计时器
type Timer struct {
	current uint32
	max 	uint32
}

// 构造函数
func NewTimer(max uint32) Timer {
	return Timer{
		current:	0,
		max:		max,
	}
}

// 重置时间
func (timer *Timer)Reset() {
	timer.current = 0
}

// 计时
func (timer *Timer)Tick() bool {
	timer.current++
	if timer.current >= timer.max {
		return true
	}
	return false
}