package server

import (
	"strings"

	"github.com/apex/log"
	"github.com/tidwall/redcon"
	"github.com/tidwall/sjson"
)

func NewResp(respAddr string, engine *Engine) *Resp {
	return &Resp{addr:respAddr, engine: engine}
}

type Resp struct {
	addr string
	engine *Engine
}

func (r *Resp) Start() {
	var ps redcon.PubSub
	var err error
	go log.WithField("addr", r.addr).Info("started resp server")
	err = redcon.ListenAndServe(r.addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			switch strings.ToLower(string(cmd.Args[0])) {
			default:
				conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
			case "ping":
				conn.WriteString("PONG")
			case "quit":
				conn.WriteString("OK")
				err := conn.Close()
				if err != nil {
					log.WithError(err).Info("error closing connection")
				} else {
					log.Info("closed connection")
				}
			case "get":
				if len(cmd.Args) != 2 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				name := string(cmd.Args[1])
				log.WithField("name", name).Info("resp get")
				info := ""
				if r.engine.HasCounter(name) {
					info, err = sjson.Set(info, "counter", r.engine.Counter(name).Count)
					if err != nil {
						log.WithError(err).Info("could not write counter")
						conn.WriteError("could not set counter")
						return
					}
				}
				if r.engine.HasTimer(name) {
					info, err = sjson.Set(info, "timer", r.engine.Timer(name).Time)
					if err != nil {
						log.WithError(err).Info("could not write timer")
						conn.WriteError("could not set timer")
						return
					}
				}
				if r.engine.HasGauge(name) {
					info, err = sjson.Set(info, "gauge", r.engine.Gauge(name).Value)
					if err != nil {
						log.WithError(err).Info("could not write gauge")
						conn.WriteError("could not set gauge")
						return
					}
				}
				if r.engine.HasSet(name) {
					info, err = sjson.Set(info, "timer", r.engine.Set(name).Buckets)
					if err != nil {
						log.WithError(err).Info("could not write timer")
						conn.WriteError("could not set timer")
						return
					}
				}
				conn.WriteBulk([]byte(info))
				log.Info("resp get done")
			case "publish":
				if len(cmd.Args) != 3 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				conn.WriteInt(ps.Publish(string(cmd.Args[1]), string(cmd.Args[2])))
			case "subscribe", "psubscribe":
				if len(cmd.Args) < 2 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				command := strings.ToLower(string(cmd.Args[0]))
				for i := 1; i < len(cmd.Args); i++ {
					if command == "psubscribe" {
						ps.Psubscribe(conn, string(cmd.Args[i]))
					} else {
						ps.Subscribe(conn, string(cmd.Args[i]))
					}
				}
			}
		},
		func(conn redcon.Conn) bool {
			// Use this function to accept or deny the connection.
			// log.Printf("accept: %s", conn.RemoteAddr())
			log.WithField("remote", conn.RemoteAddr()).Info("connected resp")
			return true
		},
		func(conn redcon.Conn, err error) {
			// This is called when the connection has been closed
			// log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
			l := log.WithField("remote", conn.RemoteAddr())
			if err != nil {
				l = l.WithError(err)
			}
			l.Info("closed resp")
		},
	)
	if err != nil {
		log.WithError(err).Warn("unexpected error with resp server")
	}
}
