package server

import (
	"io"
	"net"
	"strings"

	"github.com/apex/log"
)

func New() *Engine {
	return &Engine{}
}

// Engine implements the entire statsd server.
type Engine struct {
	conn *net.UDPConn
}

// Start begins accepting UDP stats packets from the network.
// Start blocks until the server is stopped.
func (e *Engine) Start(udpAddr, tcpAddr, httpAddr, respAddr string) error {
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

		l.WithFields(log.Fields{"sampling": sampling, "name":name, "value":value, "float": isFloat}).Info("parts")

		switch kind {
		case "c":
		case "ms":
		case "g":
		case "s":
		default:
			l.WithField("type", kind).Info("unknown type")
		}
	}
}
