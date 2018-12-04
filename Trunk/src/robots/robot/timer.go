package robot

import (
	"time"
)

type TimeCB func(interface{})
type Timer struct {
	Name     string
	delay    time.Duration
	time     time.Time
	callback TimeCB
	args     interface{}
	count    int
}

func (t *Timer) Exec(now time.Time) bool {
	if now.Sub(t.time) >= 0 {
		t.callback(t.args)
		t.time = now.Add(t.delay)
		if t.count > 0 {
			t.count--
		}
		if t.count == 0 {
			return true
		}
	}

	return false
}
