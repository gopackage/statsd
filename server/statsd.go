package server

import (
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/gopackage/statsd/stats"
	"github.com/gopackage/statsd/web"
)

func New() *Engine {
	return &Engine{
		counters: make(map[string]*stats.Counter),
		timers:   make(map[string]*stats.Timer),
		gauges:   make(map[string]*stats.Gauge),
		sets:     make(map[string]*stats.Set),
	}
}

// Engine implements the entire statsd server.
type Engine struct {
	conn     *net.UDPConn
	counters map[string]*stats.Counter
	timers   map[string]*stats.Timer
	gauges   map[string]*stats.Gauge
	sets     map[string]*stats.Set
	resp     *Resp
}

// Start begins accepting UDP stats packets from the network.
// Start blocks until the server is stopped.
func (e *Engine) Start(udpAddr, tcpAddr, httpAddr, respAddr string) error {
	webServer := web.NewServer(httpAddr)
	go webServer.Start()
	addr, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		return err
	}
	log.WithField("addr", addr).Info("Dialing")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	log.Info("Connected")

	if len(respAddr) > 0 {
		e.resp = NewResp(respAddr, e)
		go e.resp.Start()
	}

	e.conn = conn
	buffer := make([]byte, 512)
	for {
		size, remote, err := conn.ReadFrom(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		// Process the data
		l := log.WithFields(log.Fields{"size": size, "remote": remote, "data": string(buffer[0:size])})
		l.Info("stat")
		if size < 3 {
			l.Info("empty stat")
			continue
		}
		// Split into parts
		parts := strings.Split(string(buffer[0:size]), "|")
		sampling := "1"
		switch len(parts) {
		case 2:
			// missing optional sampling rate - use the default of 1 (set above)
		case 3:
			// contains optional sampling rate
			if strings.HasPrefix(parts[2], "@") { // strip out prefix
				sampling = parts[2][1:]
			} else {
				sampling = parts[2]
			}
		default:
			l.WithField("parts", len(parts)).Info("incorrect format |")
			continue
		}
		kind := parts[1]
		parts = strings.SplitN(parts[0], ":", 2)
		if len(parts) != 2 {
			l.WithField("parts", len(parts)).Info("incorrect format :")
			continue
		}
		name := parts[0]
		value := parts[1]
		isFloat := strings.Contains(value, ".")

		l.WithFields(log.Fields{"sampling": sampling, "name": name, "value": value, "float": isFloat}).Info("parts")

		switch kind {
		case "c":
			c := e.Counter(name)
			value, err := strconv.Atoi(value)
			if err != nil {
				l.WithError(err).Info("could not parse value")
				continue
			}
			c.Add(int64(value))
		case "ms":
			t := e.Timer(name)
			value, err := strconv.Atoi(value)
			if err != nil {
				l.WithError(err).Info("could not parse value")
				continue
			}
			t.Set(int64(value))
		case "g":
			g := e.Gauge(name)
			value, err := strconv.Atoi(value)
			if err != nil {
				l.WithError(err).Info("could not parse value")
				continue
			}
			g.Set(int64(value))
		case "s":
			s := e.Set(name)
			value, err := strconv.Atoi(value)
			if err != nil {
				l.WithError(err).Info("could not parse value")
				continue
			}
			s.Add(int64(value))
		default:
			l.WithField("type", kind).Info("unknown type")
		}
	}
}

func (e *Engine) HasCounter(name string) bool {
	_, ok := e.counters[name]
	return ok
}

func (e *Engine) Counter(name string) *stats.Counter {
	c, ok := e.counters[name]
	if !ok {
		c = stats.NewCounter(name)
		e.counters[name] = c
	}
	return c
}

func (e *Engine) HasTimer(name string) bool {
	_, ok := e.timers[name]
	return ok
}

func (e *Engine) Timer(name string) *stats.Timer {
	t, ok := e.timers[name]
	if !ok {
		t = stats.NewTimer(name)
		e.timers[name] = t
	}
	return t
}

func (e *Engine) HasGauge(name string) bool {
	_, ok := e.gauges[name]
	return ok
}

func (e *Engine) Gauge(name string) *stats.Gauge {
	g, ok := e.gauges[name]
	if !ok {
		g = stats.NewGauge(name)
		e.gauges[name] = g
	}
	return g
}

func (e *Engine) HasSet(name string) bool {
	_, ok := e.sets[name]
	return ok
}

func (e *Engine) Set(name string) *stats.Set {
	s, ok := e.sets[name]
	if !ok {
		s = stats.NewSet(name)
		e.sets[name] = s
	}
	return s
}
