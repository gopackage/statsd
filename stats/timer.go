package stats

import "github.com/apex/log"

func NewTimer(name string) *Timer {
	return &Timer{Name: name, Time: 0}
}

// Timer times a statistic and must have an integer value.
type Timer struct {
	Name string // Name of the timer
	Time int64 // Time is the value of the timer
}

func (t *Timer) Set(value int64) {
	t.Time = value
	log.WithFields(log.Fields{"name":t.Name, "time":t.Time}).Info("timer set")
}
