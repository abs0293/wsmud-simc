package simulator

import "math"

const (
	TICK int = 1 // 单位:毫秒(ms)
)

type Timer struct {
	rt int
}

func (t *Timer) Update(diff int) {
	if t.rt > 0 {
		t.rt -= diff
		if t.rt < 0 {
			t.rt = 0
		}
	}
}

func (t *Timer) Start(cd int) *Timer {
	t.rt = cd
	if t.rt < 0 {
		t.rt = 0
	}
	return t
}

func (t Timer) IsDone() bool {
	return t.rt == 0
}

func (t *Timer) Done() {
	t.rt = 0
}

func (t Timer) GetRemaining() int {
	return t.rt
}

func NewTimer() *Timer {
	return &Timer{}
}

func RoundFloat64(x float64, d int) float64 {
	v := math.Pow10(d)
	return math.Trunc(x*v+0.5) / v
}

func Sec2Ms(x float64) int {
	return int(RoundFloat64(x, 3) * 1000)
}

func Ms2Sec(x int) float64 {
	return RoundFloat64(float64(x)/1000, 3)
}
