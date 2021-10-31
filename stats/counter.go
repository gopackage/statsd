package stats

import "github.com/apex/log"

func NewCounter(name string) *Counter {
	return &Counter{Name: name, Count: 0}
}

// Counter counts a statistic and must have an integer value.
type Counter struct {
	Name string // Name of the counter
	Count int64 // Count is the value of the counter
}

func (c *Counter) Add(value int64) {
	c.Count += value
	log.WithFields(log.Fields{"name":c.Name, "count":c.Count}).Info("counter added")
}
