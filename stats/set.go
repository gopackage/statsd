package stats

import "github.com/apex/log"

func NewSet(name string) *Set {
	return &Set{Name: name, Buckets: make(map[int64]int64)}
}

// Set tracks occurrences of a statistic and must have an integer value.
type Set struct {
	Name string // Name of the set
	Buckets map[int64]int64 // Buckets is the number of occurrences of each item in the set
}

func (s *Set) Add(value int64) {
	b, ok := s.Buckets[value]
	if !ok {
		b = 0
	}
	b++
	s.Buckets[value] = b
	log.WithFields(log.Fields{"name":s.Name, "value": value, "count": b}).Info("set added")
}