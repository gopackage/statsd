package statsd

import (
	"bytes"
	"io"
	"net"
	"strconv"

	"github.com/apex/log"
)

// Server implements the entire statsd server.
type Server struct {
	conn *net.UDPConn
}

// Start begins accepting UDP stats packets from the network.
// Start blocks until the server is stopped.
func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", ":9045")
	if err != nil {
		return err
	}
	log.WithField("addr", addr).Info("Dialing")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	log.Info("Connected")

	s.conn = conn
	buffer := make([]byte, 512)
	for {
		size, remote, err := conn.ReadFrom(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		// TODO process the data
		log.WithFields(log.Fields{"size": size, "remote": remote, "data": string(buffer[0:size])}).Info("stat")
		if size < 3 {
			log.Info("empty stat")
			continue
		}
		switch buffer[0] {
		case byte('#'):
			// Value
			parts := bytes.SplitN(buffer[1:size], []byte("="), 2)
			if len(parts) != 2 {
				log.Info("value format")
				continue
			}
			name := string(parts[1])
			value, err := strconv.ParseFloat(string(parts[0]), 64)
			if err != nil {
				log.WithError(err).Info("parse value")
				continue
			}
			log.WithFields(log.Fields{"name": name, "value": value}).Info("value")
		default:
			parts := bytes.SplitN(buffer[:size], []byte("="), 2)
			log.WithField("parts", parts).Info("parts")
			if len(parts) != 2 {
				log.Info("value format")
				continue
			}
			name := string(parts[1])
			count, err := strconv.ParseInt(string(parts[0]), 0, 64)
			if err != nil {
				log.WithError(err).Info("parse count")
				continue
			}
			log.WithFields(log.Fields{"name": name, "count": count}).Info("count")
		}
	}
}

// Stop terminates the server (not thread safe).
func (s *Server) Stop() {
	// Not thread safe but we don't anticipate multiple calls to Stop.
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}
