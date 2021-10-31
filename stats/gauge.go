package stats

import "github.com/apex/log"

func NewGauge(name string) *Gauge {
	return &Gauge{Name: name, Value: 0}
}

// Gauge tracks a statistic and must have an integer value.
type Gauge struct {
	Name string // Name of the gauge
	Value int64 // Value of the gauge
}

func (g *Gauge) Set(value int64) {
	g.Value = value
	log.WithFields(log.Fields{"name":g.Name, "value":g.Value}).Info("gauge set")
}