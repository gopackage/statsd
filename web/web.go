// Package web adds the web UI for displaying statistics.
package web

import (
	"embed"
	"errors"
	"io"
	"net/http"

	"github.com/apex/log"
	"golang.org/x/net/websocket"
)

//go:embed app/*
var content embed.FS

func NewServer(addr string) *Server {
	return &Server{Addr: addr}
}

type Server struct {
	Addr string // Addr that the web server listens for incoming connections
}

func (s *Server) Start() {
	http.HandleFunc("/", RedirectHome)
	http.Handle("/app/", http.FileServer(http.FS(content)))
	http.Handle("/ws", websocket.Handler(s.HandleData))

	err := http.ListenAndServe(s.Addr, nil)
	if err != nil {
		log.WithError(err).Error("error running web server")
	} else {
		log.Info("web server stopped")
	}
}

func (s *Server) HandleData(ws *websocket.Conn) {
	log.Debug("ws client connected")
	defer func(ws *websocket.Conn) {
		log.Debug("ws client closing")
		err := ws.Close()
		if err != nil {
			log.WithError(err).Error("trouble closing web socket")
		}
	}(ws)
	// Read responses loop
	for {
		buf := make([]byte, 1024)
		read, err := ws.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// Normal close
				return
			}
			log.WithError(err).Error("trouble reading web socket")
			return
		}
		if read == 0 {
			log.Debug("empty recv")
			continue
		}
		log.WithFields(log.Fields{"read": read, "data": string(buf[0:read])}).Debugf("ws client recv")
		// Currently we don't do anything with data from the client
	}
}

func RedirectHome(res http.ResponseWriter, request *http.Request) {
	res.Header().Set("Location", "/app/")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
